package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dianlight/gollama.cpp"
)

// DraftSequence represents a draft sequence for speculative decoding
type DraftSequence struct {
	Active    bool
	Drafting  bool
	Skip      bool
	Tokens    []gollama.LlamaToken
	IBatchTgt []int32
}

// SpeculativeConfig holds configuration for speculative decoding
type SpeculativeConfig struct {
	MaxDraftTokens int     // Maximum number of tokens to draft
	PSplit         float64 // Probability threshold for splitting draft branches
	Temperature    float32 // Sampling temperature
}

func main() {
	var (
		targetModel = flag.String("model", "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf", "Path to the target (main) GGUF model file")
		draftModel  = flag.String("draft-model", "", "Path to the draft (faster) GGUF model file")
		prompt      = flag.String("prompt", "The future of AI is", "Prompt text to generate from")
		nPredict    = flag.Int("n-predict", 100, "Number of tokens to predict")
		nDraft      = flag.Int("n-draft", 5, "Number of tokens to draft ahead")
		threads     = flag.Int("threads", 4, "Number of threads to use")
		ctx         = flag.Int("ctx", 2048, "Context size")
		temp        = flag.Float64("temperature", 0.1, "Sampling temperature (0.0 = greedy)")
		seed        = flag.Int64("seed", -1, "Random seed (-1 for random)")
		verbose     = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	if *targetModel == "" {
		fmt.Fprintf(os.Stderr, "Error: target model path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *draftModel == "" {
		// Use the same model for both target and draft if no draft model specified
		*draftModel = *targetModel
		fmt.Println("Note: Using the same model for both target and draft (no acceleration)")
	}

	fmt.Printf("Gollama.cpp Speculative Decoding Example %s\n", gollama.FullVersion)
	fmt.Printf("Target Model: %s\n", *targetModel)
	fmt.Printf("Draft Model: %s\n", *draftModel)
	fmt.Printf("Prompt: %s\n", *prompt)
	fmt.Printf("Max draft tokens: %d\n", *nDraft)
	fmt.Printf("Temperature: %.2f\n", *temp)
	fmt.Println()

	// Set random seed
	if *seed == -1 {
		*seed = time.Now().UnixNano()
	}
	rand.Seed(*seed)

	// Initialize the backend
	fmt.Print("Initializing backend... ")
	gollama.Backend_init()
	defer gollama.Backend_free()
	fmt.Println("done")

	// Load target model
	fmt.Print("Loading target model... ")
	targetModelParams := gollama.Model_default_params()
	targetModelParams.UseMmap = 1
	targetModelParams.UseMlock = 0

	modelTgt, err := gollama.Model_load_from_file(*targetModel, targetModelParams)
	if err != nil {
		log.Fatalf("Failed to load target model: %v", err)
	}
	defer gollama.Model_free(modelTgt)
	fmt.Println("done")

	// Load draft model
	fmt.Print("Loading draft model... ")
	draftModelParams := gollama.Model_default_params()
	draftModelParams.UseMmap = 1
	draftModelParams.UseMlock = 0

	modelDft, err := gollama.Model_load_from_file(*draftModel, draftModelParams)
	if err != nil {
		log.Fatalf("Failed to load draft model: %v", err)
	}
	defer gollama.Model_free(modelDft)
	fmt.Println("done")

	// Create target context
	fmt.Print("Creating target context... ")
	ctxParamsTgt := gollama.Context_default_params()
	ctxParamsTgt.NCtx = uint32(*ctx)
	ctxParamsTgt.NThreads = int32(*threads)
	ctxParamsTgt.NThreadsBatch = int32(*threads)
	ctxParamsTgt.Logits = 1

	ctxTgt, err := gollama.Init_from_model(modelTgt, ctxParamsTgt)
	if err != nil {
		log.Fatalf("Failed to create target context: %v", err)
	}
	defer gollama.Free(ctxTgt)
	fmt.Println("done")

	// Create draft context
	fmt.Print("Creating draft context... ")
	ctxParamsDft := gollama.Context_default_params()
	ctxParamsDft.NCtx = uint32(*ctx)
	ctxParamsDft.NThreads = int32(*threads)
	ctxParamsDft.NThreadsBatch = int32(*threads)
	ctxParamsDft.Logits = 1

	ctxDft, err := gollama.Init_from_model(modelDft, ctxParamsDft)
	if err != nil {
		log.Fatalf("Failed to create draft context: %v", err)
	}
	defer gollama.Free(ctxDft)
	fmt.Println("done")

	// Tokenize the prompt
	fmt.Print("Tokenizing prompt... ")
	tokens, err := gollama.Tokenize(modelTgt, *prompt, true, false)
	if err != nil {
		log.Fatalf("Failed to tokenize: %v", err)
	}
	fmt.Printf("done (%d tokens)\n", len(tokens))

	if *verbose {
		fmt.Printf("Prompt tokens: %v\n", tokens)
	}

	// Process the prompt with both models
	fmt.Print("Processing prompt... ")

	// Target model: process all tokens except the last one
	if len(tokens) > 1 {
		promptBatchTgt := gollama.Batch_get_one(tokens[:len(tokens)-1])
		if err := gollama.Decode(ctxTgt, promptBatchTgt); err != nil {
			log.Fatalf("Failed to decode prompt (target): %v", err)
		}
		gollama.Batch_free(promptBatchTgt)
	}

	// Target model: process the last token
	lastTokenBatchTgt := gollama.Batch_get_one(tokens[len(tokens)-1:])
	if err := gollama.Decode(ctxTgt, lastTokenBatchTgt); err != nil {
		log.Fatalf("Failed to decode last token (target): %v", err)
	}
	gollama.Batch_free(lastTokenBatchTgt)

	// Draft model: process all tokens
	promptBatchDft := gollama.Batch_get_one(tokens)
	if err := gollama.Decode(ctxDft, promptBatchDft); err != nil {
		log.Fatalf("Failed to decode prompt (draft): %v", err)
	}
	gollama.Batch_free(promptBatchDft)

	fmt.Println("done")

	// Start generation
	fmt.Printf("\nGenerated text:\n%s", *prompt)

	config := SpeculativeConfig{
		MaxDraftTokens: *nDraft,
		Temperature:    float32(*temp),
	}

	// Statistics
	totalTokens := 0
	acceptedTokens := 0
	draftedTokens := 0
	generationStart := time.Now()

	// Main speculative decoding loop
	for i := 0; i < *nPredict; i++ {
		// Phase 1: Draft tokens using the draft model
		draftTokens := draftPhase(ctxDft, modelDft, config, *verbose)
		draftedTokens += len(draftTokens)

		if len(draftTokens) == 0 {
			// If no tokens were drafted, sample directly from target
			token := sampleTargetToken(ctxTgt, config.Temperature, *verbose)
			if token == gollama.LLAMA_TOKEN_NULL {
				break
			}

			piece := gollama.Token_to_piece(modelTgt, token, false)
			fmt.Print(piece)

			// Update both contexts with the accepted token
			updateContext(ctxTgt, token)
			updateContext(ctxDft, token)

			totalTokens++
			continue
		}

		// Phase 2: Verify draft tokens with target model
		acceptedCount := verifyPhase(ctxTgt, ctxDft, modelTgt, draftTokens, config, *verbose)
		acceptedTokens += acceptedCount
		totalTokens += acceptedCount

		if acceptedCount == 0 {
			// If no draft tokens were accepted, sample from target
			token := sampleTargetToken(ctxTgt, config.Temperature, *verbose)
			if token == gollama.LLAMA_TOKEN_NULL {
				break
			}

			piece := gollama.Token_to_piece(modelTgt, token, false)
			fmt.Print(piece)

			// Update both contexts
			updateContext(ctxTgt, token)
			updateContext(ctxDft, token)

			totalTokens++
		}
	}

	generationTime := time.Since(generationStart)

	// Print statistics
	fmt.Printf("\n\nSpeculative Decoding Statistics:\n")
	fmt.Printf("Total tokens generated: %d\n", totalTokens)
	fmt.Printf("Draft tokens created: %d\n", draftedTokens)
	fmt.Printf("Draft tokens accepted: %d\n", acceptedTokens)
	if draftedTokens > 0 {
		fmt.Printf("Acceptance rate: %.2f%%\n", float64(acceptedTokens)/float64(draftedTokens)*100)
	}
	fmt.Printf("Generation time: %v\n", generationTime)
	if totalTokens > 0 {
		fmt.Printf("Tokens per second: %.2f\n", float64(totalTokens)/generationTime.Seconds())
	}
}

// draftPhase generates draft tokens using the draft model
func draftPhase(ctx gollama.LlamaContext, model gollama.LlamaModel, config SpeculativeConfig, verbose bool) []gollama.LlamaToken {
	var draftTokens []gollama.LlamaToken

	for i := 0; i < config.MaxDraftTokens; i++ {
		token := sampleTargetToken(ctx, config.Temperature, verbose)
		if token == gollama.LLAMA_TOKEN_NULL {
			break
		}

		draftTokens = append(draftTokens, token)

		// Update draft context with the drafted token
		updateContext(ctx, token)

		if verbose {
			piece := gollama.Token_to_piece(model, token, false)
			fmt.Printf("[DRAFT] Token %d: %d ('%s')\n", i, token, piece)
		}
	}

	return draftTokens
}

// verifyPhase verifies draft tokens with the target model
func verifyPhase(ctxTgt, ctxDft gollama.LlamaContext, modelTgt gollama.LlamaModel, draftTokens []gollama.LlamaToken, config SpeculativeConfig, verbose bool) int {
	acceptedCount := 0

	for i, draftToken := range draftTokens {
		// Sample from target model
		targetToken := sampleTargetToken(ctxTgt, config.Temperature, verbose)

		if verbose {
			draftPiece := gollama.Token_to_piece(modelTgt, draftToken, false)
			targetPiece := gollama.Token_to_piece(modelTgt, targetToken, false)
			fmt.Printf("[VERIFY] Draft: %d ('%s'), Target: %d ('%s')\n",
				draftToken, draftPiece, targetToken, targetPiece)
		}

		if targetToken == draftToken {
			// Accept the drafted token
			piece := gollama.Token_to_piece(modelTgt, draftToken, false)
			fmt.Print(piece)

			// Update both contexts with accepted token
			updateContext(ctxTgt, draftToken)
			acceptedCount++

			if verbose {
				fmt.Printf("[ACCEPT] Token %d accepted\n", i)
			}
		} else {
			// Reject the drafted token, output the target token instead
			piece := gollama.Token_to_piece(modelTgt, targetToken, false)
			fmt.Print(piece)

			// Update target context with target token
			updateContext(ctxTgt, targetToken)
			acceptedCount++ // Count the target token as accepted

			if verbose {
				fmt.Printf("[REJECT] Token %d rejected, using target token\n", i)
			}

			// Stop verification after first rejection
			break
		}
	}

	// Resynchronize draft context with target context
	// In a real implementation, you'd need to track the context state more carefully
	// For simplicity, we'll just continue from where we left off

	return acceptedCount
}

// sampleTargetToken samples a token from the given context
func sampleTargetToken(ctx gollama.LlamaContext, temperature float32, verbose bool) gollama.LlamaToken {
	if temperature <= 0.0 {
		// Greedy sampling - find the token with highest probability
		logits := gollama.Get_logits_ith(ctx, -1)
		if logits == nil {
			return gollama.LLAMA_TOKEN_NULL
		}

		// For simplicity, we'll use the sampler from the library
		sampler := gollama.Sampler_init_greedy()
		defer gollama.Sampler_free(sampler)

		return gollama.Sampler_sample(sampler, ctx, -1)
	} else {
		// Temperature sampling would require more complex implementation
		// For now, fall back to greedy sampling
		sampler := gollama.Sampler_init_greedy()
		defer gollama.Sampler_free(sampler)

		return gollama.Sampler_sample(sampler, ctx, -1)
	}
}

// updateContext updates the context with a new token
func updateContext(ctx gollama.LlamaContext, token gollama.LlamaToken) error {
	batch := gollama.Batch_get_one([]gollama.LlamaToken{token})
	defer gollama.Batch_free(batch)

	return gollama.Decode(ctx, batch)
}

// Helper function to check if models are compatible for speculative decoding
func checkModelCompatibility(modelTgt, modelDft gollama.LlamaModel, verbose bool) error {
	// In a real implementation, you would check:
	// - Vocabulary compatibility
	// - Special tokens (BOS, EOS, etc.)
	// - Token mappings

	if verbose {
		fmt.Println("Note: Model compatibility checking is simplified in this example")
	}

	return nil
}
