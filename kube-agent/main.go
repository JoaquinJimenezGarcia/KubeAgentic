package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	kube "github.com/JoaquinJimenezGarcia/kube-agent/internal"
)

var (
	startTime    = time.Now()
	requestCount uint64
	version      = "v0.1.0"
)

type Req struct {
	Action       string              `json:"action"`
	ResourceType string              `json:"resource_type"`
	Spec         kube.DeploymentSpec `json:"spec"`
}

func main() {
	http.HandleFunc("/context", handleContext)
	http.HandleFunc("/apply", handleDeployment)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/status", handleStatus)

	log.Printf("[INFO] MCP Kubernetes provider running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}

func handleContext(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestCount, 1)

	kubeClient, err := kube.NewClient()
	if err != nil {
		log.Fatalf("[ERROR] kube.NewClient: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, err := kubeClient.GetKubeContext()
	if err != nil {
		log.Fatalf("[ERROR] kube.GetKubeContext: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] /context request succeeded")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ctx)
}

func handleDeployment(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&requestCount, 1)

	kubeClient, err := kube.NewClient()
	if err != nil {
		log.Printf("[ERROR] Couldn't stablish session with client: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	var req Req

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ResourceType != "deployment" {
		log.Printf("[ERROR] Invalid request: %v", err)
		http.Error(w, "Invalid request", 400)
		return
	}

	if req.Action == "create" {
		err = handleApply(req, kubeClient)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Apply request accepted (stubbed)",
		})
	} else if req.Action == "delete" {
		err = handleDelete(req, kubeClient)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Delete request accepted (stubbed)",
		})
	}

}

func handleApply(req Req, kubeClient *kube.KubeClient) error {
	log.Printf("[INFO] /apply received request: %v", req)

	err := kube.ApplyDeployment(kubeClient.Clientset, req.Spec)
	if err != nil {
		log.Printf("[ERROR] Couldn't create the resource: %v", err)
		return err
	}

	log.Printf("[INFO] /apply applied resources: %v", req.Spec.Name)

	return nil
}

func handleDelete(req Req, kubeClient *kube.KubeClient) error {
	log.Printf("[INFO] /delete received request: %v", req)

	err := kube.DeleteDeployment(kubeClient.Clientset, req.Spec.Name, req.Spec.Namespace)
	if err != nil {
		log.Printf("[ERROR] Couldn't delete the resource: %v", err)
		return err
	}

	log.Printf("[INFO] /delete applied resources: %v", req.Spec.Name)

	return nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] /health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	status := map[string]interface{}{
		"status":   "running",
		"uptime":   uptime.String(),
		"version":  version,
		"requests": atomic.LoadUint64(&requestCount),
		"time":     time.Now().Format(time.RFC3339),
	}

	log.Printf("[INFO] /status check")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
