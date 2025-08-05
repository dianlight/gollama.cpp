package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"unsafe"

	"gollama"
)

// Chunk represents a text chunk with metadata and embedding
type Chunk struct {
	Filename  string               // Source filename
	FilePos   int64                // Position in original file
	TextData  string               // Original text content
	Tokens    []gollama.LlamaToken // Tokenized content
	Embedding []float32            // Text embedding vector
}

// SimilarityResult represents a chunk with its similarity score to a query
type SimilarityResult struct {
	ChunkIndex int     // Index of the chunk
	Similarity float32 // Cosine similarity score
}

// RetrievalConfig holds configuration for the retrieval system
type RetrievalConfig struct {
	ChunkSize      int    // Minimum size of each text chunk
	ChunkSeparator string // String to divide chunks by
	TopK           int    // Number of top similar chunks to return
	Verbose        bool   // Enable verbose output
}

func main() {
	var (
		modelPath      = flag.String("model", "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf", "Path to the GGUF model file (should support embeddings)")
		contextFiles   = flag.String("context-files", "", "Comma-separated list of files to embed for retrieval")
		chunkSize      = flag.Int("chunk-size", 200, "Minimum size of each text chunk to be embedded")
		chunkSeparator = flag.String("chunk-separator", "\n", "String to divide chunks by")
		topK           = flag.Int("top-k", 3, "Number of top similar chunks to return")
		threads        = flag.Int("threads", 4, "Number of threads to use")
		ctx            = flag.Int("ctx", 2048, "Context size")
		verbose        = flag.Bool("verbose", false, "Enable verbose output")
		interactive    = flag.Bool("interactive", true, "Enable interactive query mode")
		query          = flag.String("query", "", "Single query to process (non-interactive mode)")
	)
	flag.Parse()

	if *modelPath == "" {
		fmt.Fprintf(os.Stderr, "Error: model path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *contextFiles == "" {
		fmt.Fprintf(os.Stderr, "Error: context files are required\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Gollama.cpp Retrieval Example %s\n", gollama.FullVersion)
	fmt.Printf("Model: %s\n", *modelPath)
	fmt.Printf("Chunk size: %d\n", *chunkSize)
	fmt.Printf("Chunk separator: %q\n", *chunkSeparator)
	fmt.Printf("Top-K: %d\n", *topK)
	fmt.Println()

	config := RetrievalConfig{
		ChunkSize:      *chunkSize,
		ChunkSeparator: *chunkSeparator,
		TopK:           *topK,
		Verbose:        *verbose,
	}

	// Parse context files
	fileList := strings.Split(*contextFiles, ",")
	for i, file := range fileList {
		fileList[i] = strings.TrimSpace(file)
	}

	// Initialize the backend
	fmt.Print("Initializing backend... ")
	err := gollama.Backend_init()
	if err != nil {
		fmt.Printf("failed (%v)\n", err)
		fmt.Println("Attempting to download llama.cpp libraries...")

		// Try to download the library
		downloadErr := gollama.LoadLibraryWithVersion("")
		if downloadErr != nil {
			log.Fatalf("Failed to download library: %v", downloadErr)
		}

		fmt.Print("Retrying backend initialization... ")
		err = gollama.Backend_init()
		if err != nil {
			log.Fatalf("Failed to initialize backend after download: %v", err)
		}
	}
	defer gollama.Backend_free()
	fmt.Println("done")

	// Load model
	fmt.Print("Loading model... ")
	modelParams := gollama.Model_default_params()
	modelParams.UseMmap = 1
	modelParams.UseMlock = 0

	model, err := gollama.Model_load_from_file(*modelPath, modelParams)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer gollama.Model_free(model)
	fmt.Println("done")

	// Create context with embeddings enabled
	fmt.Print("Creating context... ")
	ctxParams := gollama.Context_default_params()
	if *ctx > math.MaxUint32 || *ctx < 0 {
		log.Fatalf("context size %d is out of range for uint32", *ctx)
	}
	if *threads > math.MaxInt32 || *threads < math.MinInt32 {
		log.Fatalf("threads count %d is out of range for int32", *threads)
	}
	ctxParams.NCtx = uint32(*ctx)
	ctxParams.NThreads = int32(*threads)
	ctxParams.NThreadsBatch = int32(*threads)
	ctxParams.Embeddings = 1 // Enable embeddings

	llamaCtx, err := gollama.Init_from_model(model, ctxParams)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer gollama.Free(llamaCtx)
	fmt.Println("done")

	// Get model embedding dimension
	nEmbd := gollama.Model_n_embd(model)
	if *verbose {
		fmt.Printf("Model embedding dimension: %d\n", nEmbd)
	}

	// Process context files into chunks
	fmt.Print("Processing context files into chunks... ")
	var allChunks []Chunk
	for _, filename := range fileList {
		if filename == "" {
			continue
		}
		chunks, err := chunkFile(filename, config)
		if err != nil {
			log.Printf("Warning: Failed to process file %s: %v", filename, err)
			continue
		}
		allChunks = append(allChunks, chunks...)
	}
	fmt.Printf("done (%d chunks)\n", len(allChunks))

	if len(allChunks) == 0 {
		log.Fatal("No chunks were created from the input files")
	}

	// Tokenize all chunks
	fmt.Print("Tokenizing chunks... ")
	for i := range allChunks {
		tokens, err := gollama.Tokenize(model, allChunks[i].TextData, true, false)
		if err != nil {
			log.Printf("Warning: Failed to tokenize chunk %d: %v", i, err)
			continue
		}
		allChunks[i].Tokens = tokens
	}
	fmt.Println("done")

	// Generate embeddings for all chunks
	fmt.Print("Generating embeddings for chunks... ")
	err = generateEmbeddings(llamaCtx, model, allChunks, nEmbd, *verbose)
	if err != nil {
		log.Fatalf("Failed to generate embeddings: %v", err)
	}
	fmt.Println("done")

	// Clear tokens to save memory (embeddings are stored)
	for i := range allChunks {
		allChunks[i].Tokens = nil
	}

	fmt.Printf("Retrieval system ready with %d chunks\n\n", len(allChunks))

	if *interactive {
		// Interactive query loop
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter query (or 'quit' to exit): ")
			if !scanner.Scan() {
				break
			}

			queryText := strings.TrimSpace(scanner.Text())
			if queryText == "" {
				continue
			}
			if queryText == "quit" || queryText == "exit" {
				break
			}

			processQuery(llamaCtx, model, allChunks, queryText, config, nEmbd)
			fmt.Println()
		}
	} else if *query != "" {
		// Single query mode
		processQuery(llamaCtx, model, allChunks, *query, config, nEmbd)
	} else {
		fmt.Println("No query provided and interactive mode disabled")
	}

	fmt.Println("Retrieval session complete.")
}

// chunkFile splits a file into chunks based on the configuration
func chunkFile(filename string, config RetrievalConfig) ([]Chunk, error) {
	var chunks []Chunk

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", filename, err)
	}
	defer file.Close()

	// Read entire file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", filename, err)
	}

	text := string(content)
	filePos := int64(0)
	currentChunk := Chunk{
		Filename: filename,
		FilePos:  filePos,
	}

	// Split by separator
	parts := strings.Split(text, config.ChunkSeparator)

	for i, part := range parts {
		// Add the separator back except for the last part
		if i < len(parts)-1 {
			part += config.ChunkSeparator
		}

		currentChunk.TextData += part

		// If chunk is large enough or this is the last part
		if len(currentChunk.TextData) >= config.ChunkSize || i == len(parts)-1 {
			if strings.TrimSpace(currentChunk.TextData) != "" {
				chunks = append(chunks, currentChunk)
			}

			// Start new chunk
			filePos += int64(len(currentChunk.TextData))
			currentChunk = Chunk{
				Filename: filename,
				FilePos:  filePos,
			}
		}
	}

	return chunks, nil
}

// generateEmbeddings creates embeddings for all chunks
func generateEmbeddings(ctx gollama.LlamaContext, model gollama.LlamaModel, chunks []Chunk, nEmbd int32, verbose bool) error {
	for i := range chunks {
		if len(chunks[i].Tokens) == 0 {
			continue
		}

		// Create batch for this chunk
		batch := gollama.Batch_get_one(chunks[i].Tokens)
		defer gollama.Batch_free(batch)

		// Decode to get embeddings
		err := gollama.Decode(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to decode chunk %d: %v", i, err)
		}

		// Get embeddings
		embeddingsPtr := gollama.Get_embeddings(ctx)
		if embeddingsPtr == nil {
			return fmt.Errorf("failed to get embeddings for chunk %d", i)
		}

		// Convert to Go slice and normalize
		embeddings := unsafe.Slice(embeddingsPtr, nEmbd)
		embeddingsCopy := make([]float32, nEmbd)
		copy(embeddingsCopy, embeddings)

		// L2 normalize the embedding
		normalizeEmbedding(embeddingsCopy)
		chunks[i].Embedding = embeddingsCopy

		if verbose && i%10 == 0 {
			fmt.Printf("Generated embedding for chunk %d/%d\n", i+1, len(chunks))
		}
	}

	return nil
}

// processQuery handles a single query and returns similar chunks
func processQuery(ctx gollama.LlamaContext, model gollama.LlamaModel, chunks []Chunk, queryText string, config RetrievalConfig, nEmbd int32) {
	if config.Verbose {
		fmt.Printf("Processing query: %s\n", queryText)
	}

	// Tokenize query
	queryTokens, err := gollama.Tokenize(model, queryText, true, false)
	if err != nil {
		log.Printf("Failed to tokenize query: %v", err)
		return
	}

	// Create batch for query
	queryBatch := gollama.Batch_get_one(queryTokens)
	defer gollama.Batch_free(queryBatch)

	// Decode query to get embedding
	err = gollama.Decode(ctx, queryBatch)
	if err != nil {
		log.Printf("Failed to decode query: %v", err)
		return
	}

	// Get query embedding
	queryEmbeddingPtr := gollama.Get_embeddings(ctx)
	if queryEmbeddingPtr == nil {
		log.Printf("Failed to get query embedding")
		return
	}

	// Convert to Go slice and normalize
	queryEmbedding := unsafe.Slice(queryEmbeddingPtr, nEmbd)
	queryEmbeddingCopy := make([]float32, nEmbd)
	copy(queryEmbeddingCopy, queryEmbedding)
	normalizeEmbedding(queryEmbeddingCopy)

	// Compute similarities
	var similarities []SimilarityResult
	for i, chunk := range chunks {
		if chunk.Embedding == nil {
			continue
		}
		sim := cosineSimilarity(chunk.Embedding, queryEmbeddingCopy)
		similarities = append(similarities, SimilarityResult{
			ChunkIndex: i,
			Similarity: sim,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Display top-k results
	fmt.Printf("Top %d similar chunks:\n", config.TopK)
	topK := config.TopK
	if topK > len(similarities) {
		topK = len(similarities)
	}

	for i := 0; i < topK; i++ {
		result := similarities[i]
		chunk := chunks[result.ChunkIndex]

		fmt.Printf("filename: %s\n", chunk.Filename)
		fmt.Printf("filepos: %d\n", chunk.FilePos)
		fmt.Printf("similarity: %.6f\n", result.Similarity)
		fmt.Printf("textdata:\n%s\n", chunk.TextData)
		fmt.Println("--------------------")
	}
}

// normalizeEmbedding normalizes an embedding vector using L2 norm
func normalizeEmbedding(embedding []float32) {
	var sum float64 = 0
	for _, val := range embedding {
		sum += float64(val * val)
	}

	if sum > 0 {
		norm := float32(1.0 / (sum * sum)) // Simplified normalization
		for i := range embedding {
			embedding[i] *= norm
		}
	}
}

// cosineSimilarity computes cosine similarity between two normalized embedding vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
	}

	// Since vectors are normalized, cosine similarity is just the dot product
	return float32(dotProduct)
}
