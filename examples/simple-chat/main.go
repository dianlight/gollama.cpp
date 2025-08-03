package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ltarantino/gollama.cpp"
)

func main() {
	var (
		modelPath = flag.String("model", "", "Path to the GGUF model file")
		prompt    = flag.String("prompt", "The future of AI is", "Prompt text to generate from")
		nPredict  = flag.Int("n-predict", 50, "Number of tokens to predict")
		threads   = flag.Int("threads", 4, "Number of threads to use")
		ctx       = flag.Int("ctx", 2048, "Context size")
	)
	flag.Parse()

	if *modelPath == "" {
		fmt.Fprintf(os.Stderr, "Error: model path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Gollama.cpp Simple Chat Example %s\n", gollama.FullVersion)
	fmt.Printf("Model: %s\n", *modelPath)
	fmt.Printf("Prompt: %s\n", *prompt)
	fmt.Printf("Threads: %d\n", *threads)
	fmt.Printf("Context: %d\n", *ctx)
	fmt.Println()

	// Initialize the library
	fmt.Print("Initializing backend... ")
	if err := gollama.Backend_init(); err != nil {
		log.Fatalf("Failed to initialize backend: %v", err)
	}
	defer gollama.Backend_free()
	fmt.Println("done")

	// Print system information
	if gollama.Supports_gpu_offload() {
		fmt.Println("GPU offload: supported")
	} else {
		fmt.Println("GPU offload: not supported")
	}

	fmt.Printf("Memory mapping: %v\n", gollama.Supports_mmap())
	fmt.Printf("Memory locking: %v\n", gollama.Supports_mlock())
	fmt.Printf("Max devices: %d\n", gollama.Max_devices())
	fmt.Println()

	// Load model
	fmt.Print("Loading model... ")
	modelParams := gollama.Model_default_params()
	modelParams.UseMmap = 1   // true as uint8
	modelParams.UseMlock = 0  // false as uint8
	modelParams.VocabOnly = 0 // false as uint8

	model, err := gollama.Model_load_from_file(*modelPath, modelParams)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer gollama.Model_free(model)
	fmt.Println("done")

	// Create context
	fmt.Print("Creating context... ")
	ctxParams := gollama.Context_default_params()
	ctxParams.NCtx = uint32(*ctx)
	ctxParams.NBatch = 512
	ctxParams.NThreads = int32(*threads)
	ctxParams.Logits = 1 // true as uint8

	context, err := gollama.Init_from_model(model, ctxParams)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}
	defer gollama.Free(context)
	fmt.Println("done")

	// Tokenize the prompt
	fmt.Print("Tokenizing prompt... ")
	tokens, err := gollama.Tokenize(model, *prompt, true, false)
	if err != nil {
		log.Fatalf("Failed to tokenize: %v", err)
	}
	fmt.Printf("done (%d tokens)\n", len(tokens))

	// Create batch
	batch := gollama.Batch_init(int32(len(tokens)), 0, 1)
	defer gollama.Batch_free(batch)

	// Add tokens to batch
	for i, token := range tokens {
		gollama.Batch_add(batch, token, gollama.LlamaPos(i), []gollama.LlamaSeqId{0}, i == len(tokens)-1)
	}

	// Process the prompt
	fmt.Print("Processing prompt... ")
	if err := gollama.Decode(context, batch); err != nil {
		log.Fatalf("Failed to decode prompt: %v", err)
	}
	fmt.Println("done")

	// Generate tokens
	fmt.Printf("\nGenerated text:\n%s", *prompt)

	// Create sampler
	sampler := gollama.Sampler_init_greedy()
	defer gollama.Sampler_free(sampler)

	nCur := len(tokens)
	for i := 0; i < *nPredict && nCur < *ctx; i++ {
		// Get logits
		logits := gollama.Get_logits_ith(context, -1)
		if logits == nil {
			log.Fatal("Failed to get logits")
		}

		// Sample next token
		candidates := gollama.Token_data_array_init(model)
		newToken := gollama.Sampler_sample(sampler, context, candidates)

		// Convert token to text
		piece := gollama.Token_to_piece(model, newToken, false)
		fmt.Print(piece)

		// Add the token to a new batch
		batch = gollama.Batch_init(1, 0, 1)
		gollama.Batch_add(batch, newToken, gollama.LlamaPos(nCur), []gollama.LlamaSeqId{0}, true)

		// Decode the new token
		if err := gollama.Decode(context, batch); err != nil {
			log.Printf("Failed to decode token: %v", err)
			break
		}

		nCur++
	}

	fmt.Println()
	fmt.Printf("\nGenerated %d tokens.\n", *nPredict)
}
