## Model Download

Models for testing can be downloaded using the following command:

```bash
make model_download
```

This command uses the `hf.sh` script (located in `scripts/hf.sh`) to download models from Hugging Face. The script is automatically copied from the llama.cpp repository when you run `make clone-llamacpp`.

### Manual Model Download

You can also use the `hf.sh` script directly to download specific models:

```bash
# Download a specific model
./scripts/hf.sh --repo TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF --file tinyllama-1.1b-chat-v1.0.Q2_K.gguf --outdir models

# Or use a direct URL
./scripts/hf.sh https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q2_K.gguf
```

### Updating the hf.sh Script

If you need to update the `hf.sh` script to the latest version from llama.cpp:

```bash
make update-hf-script
```