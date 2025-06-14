package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func main() {
	if len(os.Args) < 2 {
		log.Print("[ERROR] Usage: kube-processor \"<your prompt here>\"")
		os.Exit(1)
	}

	userPrompt := os.Args[1]

	// Step 1: Send prompt to Ollama
	payload := OllamaRequest{
		Model: "llama3.2",
		Prompt: fmt.Sprintf(`
You are a Kubernetes assistant. Given a prompt, return only a valid JSON in this format:

{
  "action": "create",
  "resource_type": "deployment",
  "spec": {
    "name": "nginx-app",
    "namespace": "default",
    "image": "nginx:latest",
    "replicas": 2,
    "port": 80
  }
}

Prompt: "%s"
`, userPrompt),
	}

	reqBody, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[ERROR] Invalid request to ollama: %v", err)
		panic(err)
	}
	defer resp.Body.Close()

	// Streamed JSON from Ollama
	var finalJSON string
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var chunk map[string]interface{}
		decoder.Decode(&chunk)
		if s, ok := chunk["response"].(string); ok {
			finalJSON += s
		}
	}

	// Step 2: Send parsed JSON to MCP agent
	log.Printf("[INFO] 🧠 LLM Output: %v", finalJSON)

	resp2, err := http.Post("http://localhost:8080/apply", "application/json", bytes.NewBuffer([]byte(finalJSON)))
	if err != nil {
		log.Printf("[ERROR] Invalid request to agent: %v", err)
		panic(err)
	}
	body, _ := io.ReadAll(resp2.Body)

	log.Printf("[INFO] 🚀 Agent Response: %v", string(body))
}
