package main

import (
	"fmt"

	gollama "github.com/ltarantino/gollama.cpp"
)

func main() {
	params := gollama.Context_default_params()
	fmt.Printf("Default NSeqMax: %d\n", params.NSeqMax)
	fmt.Printf("Default NCtx: %d\n", params.NCtx)
	fmt.Printf("Default NBatch: %d\n", params.NBatch)
	fmt.Printf("Default NUbatch: %d\n", params.NUbatch)
}
