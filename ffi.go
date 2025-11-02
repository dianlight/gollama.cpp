package gollama

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// FFI type definitions for llama.cpp structs
var (
	// LlamaModelParams FFI type
	ffiTypeLlamaModelParams = ffi.Type{
		Type: ffi.Struct,
		Elements: &[]*ffi.Type{
			&ffi.TypePointer, // devices
			&ffi.TypePointer, // tensor_buft_overrides
			&ffi.TypeSint32,  // n_gpu_layers
			&ffi.TypeSint32,  // split_mode
			&ffi.TypeSint32,  // main_gpu
			&ffi.TypePointer, // tensor_split
			&ffi.TypePointer, // progress_callback
			&ffi.TypePointer, // progress_callback_user_data
			&ffi.TypePointer, // kv_overrides
			&ffi.TypeUint8,   // vocab_only
			&ffi.TypeUint8,   // use_mmap
			&ffi.TypeUint8,   // use_mlock
			&ffi.TypeUint8,   // check_tensors
			&ffi.TypeUint8,   // use_extra_bufts
			nil,
		}[0],
	}

	// LlamaContextParams FFI type
	ffiTypeLlamaContextParams = ffi.Type{
		Type: ffi.Struct,
		Elements: &[]*ffi.Type{
			&ffi.TypeUint32,  // seed
			&ffi.TypeUint32,  // n_ctx
			&ffi.TypeUint32,  // n_batch
			&ffi.TypeUint32,  // n_ubatch
			&ffi.TypeUint32,  // n_seq_max
			&ffi.TypeSint32,  // n_threads
			&ffi.TypeSint32,  // n_threads_batch
			&ffi.TypeSint32,  // rope_scaling_type
			&ffi.TypeSint32,  // pooling_type
			&ffi.TypeSint32,  // attention_type
			&ffi.TypeFloat,   // rope_freq_base
			&ffi.TypeFloat,   // rope_freq_scale
			&ffi.TypeFloat,   // yarn_ext_factor
			&ffi.TypeFloat,   // yarn_attn_factor
			&ffi.TypeFloat,   // yarn_beta_fast
			&ffi.TypeFloat,   // yarn_beta_slow
			&ffi.TypeUint32,  // yarn_orig_ctx
			&ffi.TypeFloat,   // defrag_thold
			&ffi.TypePointer, // cb_eval
			&ffi.TypePointer, // cb_eval_user_data
			&ffi.TypeSint32,  // type_k
			&ffi.TypeSint32,  // type_v
			&ffi.TypePointer, // abort_callback
			&ffi.TypePointer, // abort_callback_data
			&ffi.TypeUint8,   // logits
			&ffi.TypeUint8,   // embeddings
			&ffi.TypeUint8,   // offload_kqv
			&ffi.TypeUint8,   // flash_attn
			&ffi.TypeUint8,   // no_perf
			nil,
		}[0],
	}

	// LlamaSamplerChainParams FFI type
	ffiTypeLlamaSamplerChainParams = ffi.Type{
		Type: ffi.Struct,
		Elements: &[]*ffi.Type{
			&ffi.TypeUint8, // no_perf
			nil,
		}[0],
	}

	// LlamaBatch FFI type
	ffiTypeLlamaBatch = ffi.Type{
		Type: ffi.Struct,
		Elements: &[]*ffi.Type{
			&ffi.TypeSint32,  // n_tokens
			&ffi.TypePointer, // token
			&ffi.TypePointer, // embd
			&ffi.TypePointer, // pos
			&ffi.TypePointer, // n_seq_id
			&ffi.TypePointer, // seq_id
			&ffi.TypePointer, // logits
			nil,
		}[0],
	}
)

// FFI function wrappers

// ffiModelDefaultParams calls llama_model_default_params using FFI
func ffiModelDefaultParams() (LlamaModelParams, error) {
	var cif ffi.Cif
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 0, &ffiTypeLlamaModelParams); status != ffi.OK {
		return LlamaModelParams{}, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_model_default_params")
	if err != nil {
		return LlamaModelParams{}, fmt.Errorf("failed to get llama_model_default_params address: %w", err)
	}

	var result LlamaModelParams
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result))
	return result, nil
}

// ffiContextDefaultParams calls llama_context_default_params using FFI
func ffiContextDefaultParams() (LlamaContextParams, error) {
	var cif ffi.Cif
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 0, &ffiTypeLlamaContextParams); status != ffi.OK {
		return LlamaContextParams{}, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_context_default_params")
	if err != nil {
		return LlamaContextParams{}, fmt.Errorf("failed to get llama_context_default_params address: %w", err)
	}

	var result LlamaContextParams
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result))
	return result, nil
}

// ffiSamplerChainDefaultParams calls llama_sampler_chain_default_params using FFI
func ffiSamplerChainDefaultParams() (LlamaSamplerChainParams, error) {
	var cif ffi.Cif
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 0, &ffiTypeLlamaSamplerChainParams); status != ffi.OK {
		return LlamaSamplerChainParams{}, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_sampler_chain_default_params")
	if err != nil {
		return LlamaSamplerChainParams{}, fmt.Errorf("failed to get llama_sampler_chain_default_params address: %w", err)
	}

	var result LlamaSamplerChainParams
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result))
	return result, nil
}

// ffiBatchInit calls llama_batch_init using FFI
func ffiBatchInit(nTokens, embd, nSeqMax int32) (LlamaBatch, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 3, &ffiTypeLlamaBatch, aTypes...); status != ffi.OK {
		return LlamaBatch{}, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_batch_init")
	if err != nil {
		return LlamaBatch{}, fmt.Errorf("failed to get llama_batch_init address: %w", err)
	}

	var result LlamaBatch
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&nTokens),
		unsafe.Pointer(&embd),
		unsafe.Pointer(&nSeqMax),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)
	return result, nil
}

// ffiModelLoadFromFile calls llama_model_load_from_file using FFI
func ffiModelLoadFromFile(pathModel *byte, params LlamaModelParams) (LlamaModel, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypePointer, &ffiTypeLlamaModelParams}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypePointer, aTypes...); status != ffi.OK {
		return 0, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_model_load_from_file")
	if err != nil {
		return 0, fmt.Errorf("failed to get llama_model_load_from_file address: %w", err)
	}

	var result LlamaModel
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&pathModel),
		unsafe.Pointer(&params),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)

	if result == 0 {
		return 0, fmt.Errorf("failed to load model")
	}
	return result, nil
}

// ffiInitFromModel calls llama_init_from_model using FFI
func ffiInitFromModel(model LlamaModel, params LlamaContextParams) (LlamaContext, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypePointer, &ffiTypeLlamaContextParams}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypePointer, aTypes...); status != ffi.OK {
		return 0, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_init_from_model")
	if err != nil {
		return 0, fmt.Errorf("failed to get llama_init_from_model address: %w", err)
	}

	var result LlamaContext
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&model),
		unsafe.Pointer(&params),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)

	if result == 0 {
		return 0, fmt.Errorf("failed to create context")
	}
	return result, nil
}

// ffiDecode calls llama_decode using FFI
func ffiDecode(ctx LlamaContext, batch LlamaBatch) (int32, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypePointer, &ffiTypeLlamaBatch}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypeSint32, aTypes...); status != ffi.OK {
		return -1, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_decode")
	if err != nil {
		return -1, fmt.Errorf("failed to get llama_decode address: %w", err)
	}

	var result int32
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&batch),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)
	return result, nil
}

// ffiEncode calls llama_encode using FFI
func ffiEncode(ctx LlamaContext, batch LlamaBatch) (int32, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypePointer, &ffiTypeLlamaBatch}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffi.TypeSint32, aTypes...); status != ffi.OK {
		return -1, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_encode")
	if err != nil {
		return -1, fmt.Errorf("failed to get llama_encode address: %w", err)
	}

	var result int32
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&batch),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)
	return result, nil
}

// ffiBatchGetOne calls llama_batch_get_one using FFI
func ffiBatchGetOne(tokens *LlamaToken, nTokens int32) (LlamaBatch, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffi.TypePointer, &ffi.TypeSint32}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffiTypeLlamaBatch, aTypes...); status != ffi.OK {
		return LlamaBatch{}, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_batch_get_one")
	if err != nil {
		return LlamaBatch{}, fmt.Errorf("failed to get llama_batch_get_one address: %w", err)
	}

	var result LlamaBatch
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&tokens),
		unsafe.Pointer(&nTokens),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)
	return result, nil
}

// ffiSamplerChainInit calls llama_sampler_chain_init using FFI
func ffiSamplerChainInit(params LlamaSamplerChainParams) (LlamaSampler, error) {
	var cif ffi.Cif
	aTypes := []*ffi.Type{&ffiTypeLlamaSamplerChainParams}
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypePointer, aTypes...); status != ffi.OK {
		return 0, fmt.Errorf("ffi.PrepCif failed: %s", status.String())
	}

	fnAddr, err := getProcAddressPlatform(libHandle, "llama_sampler_chain_init")
	if err != nil {
		return 0, fmt.Errorf("failed to get llama_sampler_chain_init address: %w", err)
	}

	var result LlamaSampler
	aValues := []unsafe.Pointer{
		unsafe.Pointer(&params),
	}
	ffi.Call(&cif, fnAddr, unsafe.Pointer(&result), aValues...)
	return result, nil
}
