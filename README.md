# MCP Kubernetes Provider with Ollama LLM Integration

![Go Version](https://img.shields.io/badge/go-1.24.2-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-experimental-orange.svg)

This project implements a **Model Context Protocol (MCP)** provider that receives high-level user requests and applies Kubernetes configurations using natural language and a local LLM powered by [Ollama](https://ollama.com/).

> Example: The user sends a prompt like _"Create a deployment for nginx with 2 replicas"_ — the agent processes this through an LLM, generates the Kubernetes manifest, and applies it automatically to the cluster.

---

## ✨ Features

- ✅ Accepts structured MCP-style JSON requests
- 🧠 Translates natural language into Kubernetes YAML via Ollama
- 🚀 Applies manifests to a live Kubernetes cluster
- 🔄 REST API with `/apply`, `/health`, `/status` endpoints
- 🪵 Structured logging using the standard `log` package
- 🧩 Modular design: swap in different LLMs or Kubernetes backends

---

## 📦 Requirements

- Go 1.24.2+
- Access to a Kubernetes cluster (`~/.kube/config`)
- [Ollama](https://ollama.com) running locally on port `11434` with a model like `llama3.2`

---

## 🔧 Installation

```bash
git clone https://github.com/JoaquinJimenezGarcia/kubematic
cd kubematic

# In one terminal
cd kube-agent
go mod tidy
go run main.go

# In other terminal
cd kube-processor
go mod tidy
go run main.go "your request here"
```

The agent will:

1. Convert this to a natural language prompt.
2. Send the prompt to the local LLM (via Ollama).
3. Parse the response into a Kubernetes manifest.
4. Apply it to the Kubernetes cluster using the Go client.
---

## 🧠 Make sure Ollama is up and running
```bash
ollama serve
ollama pull llama3.2
```

---

## 🔌 API Endpoints
| Method | Path      | Description                   |
| ------ | --------- | ----------------------------- |
| POST   | /apply    | Accepts an MCP JSON request   |
| GET    | /context  | Reads context of Kubernetes   |
| GET    | /health   | Health check endpoint         |
| GET    | /status   | Returns internal agent status |

---

## 📝 Example Request
This is what actually the agent actually is expecting POST /apply (request can be done directly without the LLM processor):
```json
{
  "action": "create",
  "resource": "deployment",
  "params": {
    "name": "nginx-deployment",
    "namespace": "default",
    "image": "nginx",
    "replicas": 2,
    "port": 80
  }
}
```
---

## 🛠 Roadmap
- ⬜ Add support for Services, Ingress, ConfigMaps, and StatefulSets
- ⬜ Improve LLM parsing robustness
- ⬜ Add schema validation before applying manifests
- ⬜ Support for response streaming from Ollama
- ⬜ Add basic authentication to the HTTP API
- ⬜ Unit and integration tests

---
## 📄 License
MIT © 2025 — Joaquin Jimenez Garcia CloudArch