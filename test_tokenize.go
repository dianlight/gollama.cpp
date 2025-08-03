package main

import (
	"fmt"
	"log"
	"os"

	gollama "github.com/ltarantino/gollama.cpp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_tokenize.go <model_path>")
	}

	modelPath := os.Args[1]

	// Load the model
	params := gollama.DefaultModelParams()
	params.NGpuLayers = 99 // Use all layers on GPU

	model, err := gollama.LoadModelFromFile(modelPath, params)
	if err != nil {
		log.Fatal("Failed to load model:", err)
	}
	defer gollama.FreeModel(model)

	fmt.Println("Model loaded successfully")

	// Test getting vocab
	fmt.Println("Getting vocab from model...")
	// vocab := llamaModelGetVocab(model)
	// fmt.Printf("Vocab pointer: %v\n", vocab)

	// Test simple tokenization
	fmt.Println("Testing tokenization...")
	tokens, err := gollama.Tokenize(model, "Hello world", false, false)
	if err != nil {
		log.Fatal("Failed to tokenize:", err)
	}

	fmt.Printf("Tokens: %v\n", tokens)
	fmt.Printf("Number of tokens: %d\n", len(tokens))
}
