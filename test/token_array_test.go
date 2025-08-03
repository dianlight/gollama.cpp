package gollama

import (
	"testing"
	"unsafe"
)

func TestTokenDataArrayFromLogits(t *testing.T) {
	// Create dummy logits array
	logits := make([]float32, 256)
	for i := 0; i < 256; i++ {
		logits[i] = float32(i) * 0.1
	}

	// Call our function with the logits
	// We don't need a real model since the function doesn't use it currently
	tokenArray := Token_data_array_from_logits(LlamaModel(0), &logits[0])

	if tokenArray == nil {
		t.Fatal("Token array should not be nil")
	}

	// Check that the size is now 256, not 32000
	expectedSize := uint64(256)
	if tokenArray.Size != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, tokenArray.Size)
	}

	// Check that Selected is initialized correctly
	if tokenArray.Selected != -1 {
		t.Errorf("Expected Selected to be -1, got %d", tokenArray.Selected)
	}

	// Check that Sorted is initialized correctly
	if tokenArray.Sorted != 0 {
		t.Errorf("Expected Sorted to be 0, got %d", tokenArray.Sorted)
	}

	// Check that we can access the first element
	if tokenArray.Data == nil {
		t.Fatal("Data pointer should not be nil")
	}

	firstToken := tokenArray.Data
	if firstToken.Id != 0 {
		t.Errorf("Expected first token ID to be 0, got %d", firstToken.Id)
	}

	if firstToken.Logit != 0.0 {
		t.Errorf("Expected first token logit to be 0.0, got %f", firstToken.Logit)
	}

	// Check that we can access the last element (index 255)
	lastElement := (*LlamaTokenData)(unsafe.Pointer(uintptr(unsafe.Pointer(tokenArray.Data)) + uintptr(255)*unsafe.Sizeof(LlamaTokenData{})))
	expectedId := LlamaToken(255)
	if lastElement.Id != expectedId {
		t.Errorf("Expected last token ID to be %d, got %d", expectedId, lastElement.Id)
	}

	expectedLogit := float32(255) * 0.1
	if lastElement.Logit != expectedLogit {
		t.Errorf("Expected last token logit to be %f, got %f", expectedLogit, lastElement.Logit)
	}

	t.Logf("SUCCESS: Token array created with size %d", tokenArray.Size)
	t.Logf("Data pointer: %p", tokenArray.Data)
	t.Logf("First token: ID=%d, Logit=%f", firstToken.Id, firstToken.Logit)
	t.Logf("Last token: ID=%d, Logit=%f", lastElement.Id, lastElement.Logit)
}
