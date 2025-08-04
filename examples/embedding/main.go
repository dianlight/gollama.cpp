package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"unsafe"

	"github.com/dianlight/gollama.cpp"
)

// splitLines splits a string into lines based on a separator
func splitLines(s, separator string) []string {
	if separator == "" {
		separator = "\n"
	}
	return strings.Split(s, separator)
}

// addSequenceToBatch adds a sequence of tokens to a batch
func addSequenceToBatch(batch *gollama.LlamaBatch, tokens []gollama.LlamaToken, seqId gollama.LlamaSeqId) {
	for i, token := range tokens {
		// We need to manually populate the batch since there's no direct helper
		// This is a simplified version - in a real implementation you'd want proper batch management
		if int32(i) < batch.NTokens {
			// Access batch data directly (unsafe but necessary for this example)
			tokensPtr := (*[1 << 20]gollama.LlamaToken)(unsafe.Pointer(batch.Token))
			posPtr := (*[1 << 20]gollama.LlamaPos)(unsafe.Pointer(batch.Pos))
			seqIdPtr := (*[1 << 20]*gollama.LlamaSeqId)(unsafe.Pointer(batch.SeqId))
			logitsPtr := (*[1 << 20]int8)(unsafe.Pointer(batch.Logits))

			tokensPtr[i] = token
			posPtr[i] = gollama.LlamaPos(i)
			// Set sequence ID (simplified)
			seqIdPtr[i] = &seqId
			logitsPtr[i] = 1 // Enable logits for last token
		}
	}
	batch.NTokens = int32(len(tokens))
}

// normalizeEmbedding normalizes an embedding vector using L2 norm (Euclidean)
func normalizeEmbedding(embedding []float32) {
	var sum float64 = 0
	for _, val := range embedding {
		sum += float64(val * val)
	}
	norm := math.Sqrt(sum)
	if norm > 0 {
		for i := range embedding {
			embedding[i] = float32(float64(embedding[i]) / norm)
		}
	}
}

// cosineSimilarity computes cosine similarity between two embedding vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}

func main() {
	var (
		modelPath = flag.String("model", "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf", "Path to the GGUF model file")
		prompt    = flag.String("prompt", "Hello World!", "Text to generate embeddings for (use | to separate multiple texts)")
		separator = flag.String("separator", "|", "Separator for multiple prompts")
		normalize = flag.Bool("normalize", true, "Normalize embeddings using L2 norm")
		threads   = flag.Int("threads", 4, "Number of threads to use")
		ctx       = flag.Int("ctx", 2048, "Context size")
		verbose   = flag.Bool("verbose", false, "Verbose output")
		outputFmt = flag.String("output-format", "default", "Output format: default, json, array")
	)
	flag.Parse()

	if *modelPath == "" {
		fmt.Fprintf(os.Stderr, "Error: model path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Gollama.cpp Embedding Example %s\n", gollama.FullVersion)
	fmt.Printf("Model: %s\n", *modelPath)

	// Initialize the backend
	gollama.Backend_init()
	defer gollama.Backend_free()

	// Load model
	modelParams := gollama.Model_default_params()
	model, err := gollama.Model_load_from_file(*modelPath, modelParams)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer gollama.Model_free(model)

	// Create context with embeddings enabled
	ctxParams := gollama.Context_default_params()
	ctxParams.NCtx = uint32(*ctx)
	ctxParams.NThreads = int32(*threads)
	ctxParams.NThreadsBatch = int32(*threads)
	ctxParams.Embeddings = 1 // Enable embeddings

	llamaCtx, err := gollama.Init_from_model(model, ctxParams)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer gollama.Free(llamaCtx)

	// Get model information
	nEmbd := gollama.Model_n_embd(model)
	if *verbose {
		fmt.Printf("Model embedding dimension: %d\n", nEmbd)
	}

	// Split prompts
	prompts := splitLines(*prompt, *separator)
	fmt.Printf("Processing %d prompt(s)\n", len(prompts))

	// Store all embeddings
	allEmbeddings := make([][]float32, len(prompts))

	// Process each prompt
	for i, promptText := range prompts {
		promptText = strings.TrimSpace(promptText)
		if promptText == "" {
			continue
		}

		if *verbose {
			fmt.Printf("Processing prompt %d: '%s'\n", i+1, promptText)
		}

		// Tokenize the prompt
		tokens, err := gollama.Tokenize(model, promptText, true, true)
		if err != nil {
			log.Printf("Failed to tokenize prompt %d: %v", i+1, err)
			continue
		}
		if len(tokens) == 0 {
			fmt.Printf("Warning: empty tokenization for prompt %d\n", i+1)
			continue
		}

		if *verbose {
			fmt.Printf("Tokenized to %d tokens\n", len(tokens))
		}

		// Create batch
		batch := gollama.Batch_init(int32(len(tokens)), 0, 1)
		defer gollama.Batch_free(batch)

		// Add tokens to batch
		addSequenceToBatch(&batch, tokens, gollama.LlamaSeqId(i))

		// Decode to get embeddings
		err = gollama.Decode(llamaCtx, batch)
		if err != nil {
			log.Printf("Failed to decode batch for prompt %d: %v", i+1, err)
			continue
		}

		// Get embeddings
		embeddingsPtr := gollama.Get_embeddings(llamaCtx)
		if embeddingsPtr == nil {
			log.Printf("Failed to get embeddings for prompt %d", i+1)
			continue
		}

		// Convert to Go slice
		embeddings := unsafe.Slice(embeddingsPtr, nEmbd)
		embeddingsCopy := make([]float32, nEmbd)
		copy(embeddingsCopy, embeddings)

		// Normalize if requested
		if *normalize {
			normalizeEmbedding(embeddingsCopy)
		}

		allEmbeddings[i] = embeddingsCopy

		// Output individual embedding
		switch *outputFmt {
		case "json":
			fmt.Printf("{\n  \"prompt\": \"%s\",\n  \"embedding\": [", promptText)
			for j, val := range embeddingsCopy {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%.6f", val)
			}
			fmt.Println("]\n}")
		case "array":
			fmt.Print("[")
			for j, val := range embeddingsCopy {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%.6f", val)
			}
			fmt.Println("]")
		default:
			fmt.Printf("Embedding %d: ", i+1)
			// Show first 5 and last 5 dimensions for readability
			if nEmbd <= 10 {
				for _, val := range embeddingsCopy {
					fmt.Printf("%.6f ", val)
				}
			} else {
				for j := 0; j < 5; j++ {
					fmt.Printf("%.6f ", embeddingsCopy[j])
				}
				fmt.Print("... ")
				for j := nEmbd - 5; j < nEmbd; j++ {
					fmt.Printf("%.6f ", embeddingsCopy[j])
				}
			}
			fmt.Println()
		}
	}

	// Compute and display cosine similarity matrix if multiple prompts
	if len(prompts) > 1 && len(allEmbeddings) > 1 && *outputFmt == "default" {
		fmt.Println("\nCosine Similarity Matrix:")
		fmt.Println()

		// Print header
		for i, prompt := range prompts {
			if i < len(allEmbeddings) && allEmbeddings[i] != nil {
				fmt.Printf("%8.8s ", prompt)
			}
		}
		fmt.Println()

		// Print similarity matrix
		for i, embA := range allEmbeddings {
			if embA == nil {
				continue
			}
			for _, embB := range allEmbeddings {
				if embB == nil {
					fmt.Printf("%8s ", "N/A")
					continue
				}
				sim := cosineSimilarity(embA, embB)
				fmt.Printf("%8.3f ", sim)
			}
			if i < len(prompts) {
				fmt.Printf(" %s", prompts[i])
			}
			fmt.Println()
		}
	}

	fmt.Println("\nEmbedding generation complete!")
}
