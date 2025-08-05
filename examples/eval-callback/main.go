package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/dianlight/gollama.cpp"
)

// EvalCallbackData simulates the callback data structure from the C++ example
type EvalCallbackData struct {
	OperationCount  int
	TensorCount     int
	BytesProcessed  int64
	StartTime       time.Time
	LastOpTime      time.Time
	EnableLogging   bool
	EnableProgress  bool
	PrintTensorData bool
	MaxLoggedOps    int
}

// NewEvalCallbackData creates a new callback data structure
func NewEvalCallbackData(enableLogging, enableProgress, printTensorData bool, maxLoggedOps int) *EvalCallbackData {
	return &EvalCallbackData{
		StartTime:       time.Now(),
		LastOpTime:      time.Now(),
		EnableLogging:   enableLogging,
		EnableProgress:  enableProgress,
		PrintTensorData: printTensorData,
		MaxLoggedOps:    maxLoggedOps,
	}
}

// SimulatedTensorInfo represents tensor information that would be available in a real callback
type SimulatedTensorInfo struct {
	Name       string
	Type       string
	Operation  string
	Dimensions []int64
	SizeBytes  int64
	IsHost     bool
	DataType   string
}

// logOperation simulates what the ggml_debug callback would log
func (cb *EvalCallbackData) logOperation(tensor SimulatedTensorInfo, srcTensors []SimulatedTensorInfo) {
	if !cb.EnableLogging {
		return
	}

	cb.OperationCount++
	cb.TensorCount++
	cb.BytesProcessed += tensor.SizeBytes

	// Limit logging to prevent overwhelming output
	if cb.MaxLoggedOps > 0 && cb.OperationCount > cb.MaxLoggedOps {
		if cb.OperationCount == cb.MaxLoggedOps+1 {
			fmt.Printf("                              ... (limiting output to %d operations)\n", cb.MaxLoggedOps)
		}
		return
	}

	currentTime := time.Now()
	timeSinceStart := currentTime.Sub(cb.StartTime)
	timeSinceLastOp := currentTime.Sub(cb.LastOpTime)
	cb.LastOpTime = currentTime

	// Format source tensor info
	var srcStr strings.Builder
	for i, src := range srcTensors {
		if i > 0 {
			srcStr.WriteString(", ")
		}
		srcStr.WriteString(fmt.Sprintf("%s{%s}", src.Name, formatDimensions(src.Dimensions)))
	}

	// Main operation log (similar to C++ eval-callback output)
	fmt.Printf("ggml_debug: %24s = (%s) %10s(%s) = {%s}\n",
		tensor.Name,
		tensor.DataType,
		tensor.Operation,
		srcStr.String(),
		formatDimensions(tensor.Dimensions))

	// Additional timing and size information
	fmt.Printf("                              └─ Op #%d, %d bytes, %s memory, %.3fms since start, %.3fms since last\n",
		cb.OperationCount,
		tensor.SizeBytes,
		getMemoryLocation(tensor.IsHost),
		timeSinceStart.Seconds()*1000,
		timeSinceLastOp.Seconds()*1000)

	// Simulate tensor data printing (like the C++ example does for non-quantized tensors)
	if cb.PrintTensorData && !strings.Contains(tensor.DataType, "q") {
		cb.printSimulatedTensorData(tensor)
	}

	fmt.Println()
}

// printSimulatedTensorData simulates printing a few tensor values
func (cb *EvalCallbackData) printSimulatedTensorData(tensor SimulatedTensorInfo) {
	if len(tensor.Dimensions) == 0 {
		return
	}

	fmt.Printf("                              Data preview (first few values):\n")
	fmt.Printf("                                     [\n")

	// Simulate printing first few values of the tensor
	rows := min(3, int(tensor.Dimensions[0]))
	cols := min(8, int(tensor.Dimensions[len(tensor.Dimensions)-1]))

	for i := 0; i < rows; i++ {
		fmt.Printf("                                       [")
		for j := 0; j < cols; j++ {
			// Generate some fake values for demonstration
			value := float32(i*10+j) * 0.123
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%8.4f", value)
		}
		if cols < int(tensor.Dimensions[len(tensor.Dimensions)-1]) {
			fmt.Printf(", ...")
		}
		fmt.Printf("]\n")
	}
	if rows < int(tensor.Dimensions[0]) {
		fmt.Printf("                                       ...\n")
	}
	fmt.Printf("                                     ]\n")
}

// showProgress displays progress information
func (cb *EvalCallbackData) showProgress() {
	if !cb.EnableProgress {
		return
	}

	elapsed := time.Since(cb.StartTime)
	avgOpsPerSec := float64(cb.OperationCount) / elapsed.Seconds()
	avgBytesPerSec := float64(cb.BytesProcessed) / elapsed.Seconds()

	fmt.Printf("\nProgress Update:\n")
	fmt.Printf("  Operations processed: %d\n", cb.OperationCount)
	fmt.Printf("  Tensors processed: %d\n", cb.TensorCount)
	fmt.Printf("  Data processed: %s\n", formatBytes(cb.BytesProcessed))
	fmt.Printf("  Elapsed time: %.2fs\n", elapsed.Seconds())
	fmt.Printf("  Average ops/sec: %.1f\n", avgOpsPerSec)
	fmt.Printf("  Average throughput: %s/sec\n", formatBytes(int64(avgBytesPerSec)))
	fmt.Println()
}

// simulateInferenceWithCallbacks demonstrates what eval callbacks would show during inference
func simulateInferenceWithCallbacks(cb *EvalCallbackData, tokens []gollama.LlamaToken, model gollama.LlamaModel) {
	fmt.Printf("=== Starting Evaluation with Callbacks ===\n")
	fmt.Printf("Tokens to process: %d\n", len(tokens))
	fmt.Printf("Logging enabled: %v\n", cb.EnableLogging)
	fmt.Printf("Progress enabled: %v\n", cb.EnableProgress)
	fmt.Printf("Tensor data printing: %v\n", cb.PrintTensorData)
	fmt.Println()

	// Simulate the operations that would occur during model inference
	// These are based on typical transformer operations

	layerCount := 22 // TinyLlama has 22 layers
	seqLen := len(tokens)

	for layer := 0; layer < layerCount; layer++ {
		// Simulate attention operations
		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("attn_q_layer_%d", layer),
			Type:       "attn_q",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), 2048},
			SizeBytes:  int64(seqLen * 2048 * 4), // float32
			IsHost:     false,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("inp_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("wq_%d", layer), Dimensions: []int64{2048, 2048}},
		})

		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("attn_k_layer_%d", layer),
			Type:       "attn_k",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), 2048},
			SizeBytes:  int64(seqLen * 2048 * 4),
			IsHost:     false,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("inp_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("wk_%d", layer), Dimensions: []int64{2048, 2048}},
		})

		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("attn_v_layer_%d", layer),
			Type:       "attn_v",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), 2048},
			SizeBytes:  int64(seqLen * 2048 * 4),
			IsHost:     false,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("inp_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("wv_%d", layer), Dimensions: []int64{2048, 2048}},
		})

		// Simulate attention scaling and softmax
		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("attn_scores_layer_%d", layer),
			Type:       "attn_scores",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), int64(seqLen)},
			SizeBytes:  int64(seqLen * seqLen * 4),
			IsHost:     true,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("attn_q_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("attn_k_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
		})

		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("attn_soft_layer_%d", layer),
			Type:       "attn_soft",
			Operation:  "SOFT_MAX",
			Dimensions: []int64{int64(seqLen), int64(seqLen)},
			SizeBytes:  int64(seqLen * seqLen * 4),
			IsHost:     true,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("attn_scores_layer_%d", layer), Dimensions: []int64{int64(seqLen), int64(seqLen)}},
		})

		// Simulate FFN operations
		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("ffn_gate_layer_%d", layer),
			Type:       "ffn_gate",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), 5632},
			SizeBytes:  int64(seqLen * 5632 * 4),
			IsHost:     false,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("attn_out_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("w_gate_%d", layer), Dimensions: []int64{2048, 5632}},
		})

		cb.logOperation(SimulatedTensorInfo{
			Name:       fmt.Sprintf("ffn_up_layer_%d", layer),
			Type:       "ffn_up",
			Operation:  "MUL_MAT",
			Dimensions: []int64{int64(seqLen), 5632},
			SizeBytes:  int64(seqLen * 5632 * 4),
			IsHost:     false,
			DataType:   "f32",
		}, []SimulatedTensorInfo{
			{Name: fmt.Sprintf("attn_out_layer_%d", layer), Dimensions: []int64{int64(seqLen), 2048}},
			{Name: fmt.Sprintf("w_up_%d", layer), Dimensions: []int64{2048, 5632}},
		})

		// Show progress every few layers
		if cb.EnableProgress && (layer%5 == 0 || layer == layerCount-1) {
			cb.showProgress()
		}

		// Add small delay to make the simulation more realistic
		time.Sleep(10 * time.Millisecond)
	}

	// Final output layer
	cb.logOperation(SimulatedTensorInfo{
		Name:       "output_logits",
		Type:       "output",
		Operation:  "MUL_MAT",
		Dimensions: []int64{int64(seqLen), 32000}, // vocab size
		SizeBytes:  int64(seqLen * 32000 * 4),
		IsHost:     true,
		DataType:   "f32",
	}, []SimulatedTensorInfo{
		{Name: "final_norm", Dimensions: []int64{int64(seqLen), 2048}},
		{Name: "output_weight", Dimensions: []int64{2048, 32000}},
	})

	fmt.Printf("=== Evaluation Complete ===\n")
	cb.showProgress()
}

// Helper functions
func formatDimensions(dims []int64) string {
	var parts []string
	for _, dim := range dims {
		parts = append(parts, fmt.Sprintf("%d", dim))
	}
	return strings.Join(parts, ", ")
}

func getMemoryLocation(isHost bool) string {
	if isHost {
		return "CPU"
	}
	return "GPU"
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var (
		modelPath       = flag.String("model", "../../models/tinyllama-1.1b-chat-v1.0.Q2_K.gguf", "Path to the GGUF model file")
		prompt          = flag.String("prompt", "The future of AI is", "Prompt text to evaluate")
		threads         = flag.Int("threads", 4, "Number of threads to use")
		ctx             = flag.Int("ctx", 512, "Context size")
		enableLogging   = flag.Bool("enable-logging", true, "Enable operation logging (simulates ggml_debug callback)")
		enableProgress  = flag.Bool("enable-progress", true, "Enable progress updates")
		printTensorData = flag.Bool("print-tensor-data", false, "Print tensor data values (verbose)")
		maxLoggedOps    = flag.Int("max-logged-ops", 50, "Maximum number of operations to log (0 = unlimited)")
		simulateOnly    = flag.Bool("simulate-only", false, "Only run simulation without actual model inference")
	)
	flag.Parse()

	if *modelPath == "" {
		fmt.Fprintf(os.Stderr, "Error: model path is required\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Gollama.cpp Evaluation Callback Example %s\n", gollama.FullVersion)
	fmt.Printf("Model: %s\n", *modelPath)
	fmt.Printf("Prompt: %s\n", *prompt)
	fmt.Printf("Threads: %d\n", *threads)
	fmt.Printf("Context: %d\n", *ctx)
	fmt.Printf("Simulation only: %v\n", *simulateOnly)
	fmt.Println()

	if *simulateOnly {
		fmt.Println("Running in simulation-only mode...")
		fmt.Println("This demonstrates what eval callbacks would show during inference.")
		fmt.Println()

		// Create callback data
		cbData := NewEvalCallbackData(*enableLogging, *enableProgress, *printTensorData, *maxLoggedOps)

		// Simulate tokenized prompt (fake token IDs)
		tokens := make([]gollama.LlamaToken, len(strings.Fields(*prompt))+2) // +2 for BOS/EOS
		for i := range tokens {
			if 1000+i > math.MaxInt32 {
				log.Fatalf("token ID %d is out of range for LlamaToken", 1000+i)
			}
			tokens[i] = gollama.LlamaToken(1000 + i) // fake token IDs
		}

		simulateInferenceWithCallbacks(cbData, tokens, 0)
		return
	}

	// Initialize the library
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

	// Print system information
	if gollama.Supports_gpu_offload() {
		fmt.Println("GPU offload: supported")
	} else {
		fmt.Println("GPU offload: not supported")
	}

	// Load model
	fmt.Print("Loading model... ")
	modelParams := gollama.Model_default_params()
	modelParams.UseMmap = 1
	modelParams.UseMlock = 0
	modelParams.VocabOnly = 0

	model, err := gollama.Model_load_from_file(*modelPath, modelParams)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer gollama.Model_free(model)
	fmt.Println("done")

	// Create context
	fmt.Print("Creating context... ")
	ctxParams := gollama.Context_default_params()
	if *ctx > math.MaxUint32 || *ctx < 0 {
		log.Fatalf("context size %d is out of range for uint32", *ctx)
	}
	if *threads > math.MaxInt32 || *threads < math.MinInt32 {
		log.Fatalf("threads count %d is out of range for int32", *threads)
	}
	ctxParams.NCtx = uint32(*ctx)
	ctxParams.NBatch = 512
	ctxParams.NSeqMax = 1
	ctxParams.NThreads = int32(*threads)
	ctxParams.Logits = 1

	// NOTE: In a real implementation, we would set eval callbacks here:
	// ctxParams.CbEval = callbackFunctionPointer
	// ctxParams.CbEvalUserData = callbackDataPointer
	// However, this requires unsafe pointer manipulation and CGO

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
		log.Fatalf("Failed to tokenize prompt: %v", err)
	}
	fmt.Printf("done (%d tokens)\n", len(tokens))

	// Create callback data for demonstration
	cbData := NewEvalCallbackData(*enableLogging, *enableProgress, *printTensorData, *maxLoggedOps)

	fmt.Println("\n=== Actual Model Inference ===")
	fmt.Println("Note: Real eval callbacks would require C callback implementation.")
	fmt.Println("This example shows what information would be available:")
	fmt.Println()

	// Run the simulation alongside actual inference
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to let inference start
		simulateInferenceWithCallbacks(cbData, tokens, model)
	}()

	// Decode tokens (this is where eval callbacks would be triggered in reality)
	fmt.Print("Evaluating tokens... ")
	startTime := time.Now()

	batch := gollama.Batch_get_one(tokens)

	if err := gollama.Decode(context, batch); err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	evalTime := time.Since(startTime)
	fmt.Printf("done (%.2fs)\n", evalTime.Seconds())

	// Get logits for the last token
	tokensLen := len(tokens)
	if tokensLen == 0 {
		log.Fatal("No tokens to get logits for")
	}
	if tokensLen-1 > math.MaxInt32 {
		log.Fatalf("token index %d is out of range for int32", tokensLen-1)
	}
	logits := gollama.Get_logits_ith(context, int32(tokensLen-1))
	if logits == nil {
		log.Fatalf("Failed to get logits")
	}

	// Print performance information
	fmt.Println("\n=== Performance Information ===")
	fmt.Printf("Evaluation time: %.2f ms\n", evalTime.Seconds()*1000)
	fmt.Printf("Tokens processed: %d\n", len(tokens))
	fmt.Printf("Processing speed: %.2f tokens/s\n", float64(len(tokens))/evalTime.Seconds())

	fmt.Println("\n=== Summary ===")
	fmt.Println("This example demonstrates:")
	fmt.Println("1. Simulated eval callbacks showing tensor operations during inference")
	fmt.Println("2. Operation counting and timing information")
	fmt.Println("3. Memory location tracking (CPU vs GPU)")
	fmt.Println("4. Tensor dimension and size monitoring")
	fmt.Println("5. Progress reporting during evaluation")
	fmt.Println()
	fmt.Println("In a real callback implementation, you would:")
	fmt.Println("- Set ctxParams.CbEval to point to your callback function")
	fmt.Println("- Set ctxParams.CbEvalUserData to point to your callback data")
	fmt.Println("- Implement the callback in C/Go using CGO for low-level access")
	fmt.Println("- Access actual tensor data, names, and operations during graph execution")

	// Wait a moment for the simulation to complete
	time.Sleep(100 * time.Millisecond)
}
