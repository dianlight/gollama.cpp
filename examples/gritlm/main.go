package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"unsafe"

	gollama "github.com/dianlight/gollama.cpp"
)

const embeddingInstruction = "<|embed|>"

// addSequenceToBatch adds a sequence of tokens to a batch
func addSequenceToBatch(batch *gollama.LlamaBatch, tokens []gollama.LlamaToken, seqId gollama.LlamaSeqId) {
	for i, token := range tokens {
		if i >= math.MaxInt32 {
			log.Fatalf("token index %d is out of range for int32", i)
		}
		if int32(i) < batch.NTokens {
			// Access batch data directly (unsafe but necessary for this example)
			tokensPtr := (*[1 << 20]gollama.LlamaToken)(unsafe.Pointer(batch.Token))
			posPtr := (*[1 << 20]gollama.LlamaPos)(unsafe.Pointer(batch.Pos))
			seqIdPtr := (*[1 << 20]*gollama.LlamaSeqId)(unsafe.Pointer(batch.SeqId))
			logitsPtr := (*[1 << 20]int8)(unsafe.Pointer(batch.Logits))

			tokensPtr[i] = token
			if i > math.MaxInt32 {
				log.Fatalf("position %d is out of range for LlamaPos", i)
			}
			posPtr[i] = gollama.LlamaPos(i)
			seqIdPtr[i] = &seqId
			// Enable outputs for all tokens in embedding mode
			logitsPtr[i] = 1
		}
	}
	tokensLen := len(tokens)
	if tokensLen > math.MaxInt32 {
		log.Fatalf("too many tokens: %d, maximum supported: %d", tokensLen, math.MaxInt32)
	}
	batch.NTokens = int32(tokensLen)
}

// normalizeEmbedding normalizes an embedding vector using L2 norm
func normalizeEmbedding(input []float32, output []float32) {
	if len(input) != len(output) {
		return
	}

	var sum float64 = 0
	for _, val := range input {
		sum += float64(val * val)
	}

	norm := math.Sqrt(sum)
	if norm > 0 {
		for i := range input {
			output[i] = float32(float64(input[i]) / norm)
		}
	} else {
		copy(output, input)
	}
}

// cosineSimilarity computes cosine similarity between two embedding vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
	}

	return float32(dotProduct)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <model-path>\n", os.Args[0])
		os.Exit(1)
	}

	modelPath := os.Args[1]

	// Initialize the backend
	gollama.Backend_init()
	defer gollama.Backend_free()

	// Load model
	fmt.Printf("Loading GritLM model from: %s\n", modelPath)
	modelParams := gollama.Model_default_params()
	modelParams.NGpuLayers = 32 // Enable GPU offloading

	model, err := gollama.Model_load_from_file(modelPath, modelParams)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer gollama.Model_free(model)
	fmt.Printf("Model loaded successfully\n")

	// Create context for embeddings
	ctxParams := gollama.Context_default_params()
	ctxParams.NCtx = 512
	ctxParams.NThreads = 4
	ctxParams.NThreadsBatch = 4
	ctxParams.Embeddings = 1 // Enable embeddings

	ctx, err := gollama.Init_from_model(model, ctxParams)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer gollama.Free(ctx)

	fmt.Printf("Context created for GritLM embeddings\n")

	// Test just one simple sentence first
	sentence := "Hello world"
	fmt.Printf("Generating embedding for: %s\n", sentence)

	// Prepare input with instruction
	inputString := embeddingInstruction + sentence

	// Tokenize the full input
	tokens, err := gollama.Tokenize(model, inputString, true, false)
	if err != nil {
		log.Fatalf("Failed to tokenize: %v", err)
	}

	tokensLen := len(tokens)
	if tokensLen > math.MaxInt32 {
		log.Fatalf("too many tokens: %d, maximum supported: %d", tokensLen, math.MaxInt32)
	}
	nToks := int32(tokensLen)
	if nToks == 0 {
		log.Fatalf("Empty tokenization")
	}

	fmt.Printf("Tokenized to %d tokens\n", nToks)

	// Create batch
	batch := gollama.Batch_init(nToks, 0, 1)

	// Add tokens to batch
	addSequenceToBatch(&batch, tokens, gollama.LlamaSeqId(0))

	// Clear previous kv_cache values (irrelevant for embeddings)
	gollama.Memory_clear(ctx, true)
	gollama.Set_causal_attn(ctx, false)

	fmt.Printf("About to decode...\n")

	// Run the model
	err = gollama.Decode(ctx, batch)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	fmt.Printf("Decode successful! Getting embeddings...\n")

	// Try standard embeddings
	embPtr := gollama.Get_embeddings(ctx)
	if embPtr == nil {
		log.Fatalf("Failed to get embeddings")
	}

	// Get embedding dimensions
	nEmbd := gollama.Model_n_embd(model)

	// Convert to Go slice
	embeddings := unsafe.Slice(embPtr, nEmbd)
	embeddingsCopy := make([]float32, nEmbd)
	copy(embeddingsCopy, embeddings)

	// Normalize the embedding (L2 norm)
	embNorm := make([]float32, nEmbd)
	normalizeEmbedding(embeddingsCopy, embNorm)

	fmt.Printf("Successfully generated embedding!\n")
	fmt.Printf("Embedding dimension: %d\n", len(embNorm))
	fmt.Printf("First 5 values: %.6f %.6f %.6f %.6f %.6f\n",
		embNorm[0], embNorm[1], embNorm[2], embNorm[3], embNorm[4])

	fmt.Println("Basic embedding generation completed successfully!")
}
