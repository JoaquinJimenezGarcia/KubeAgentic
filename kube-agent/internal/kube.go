package kube

import (
	"context"
	"fmt"

	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	Clientset *kubernetes.Clientset
}

// NewClient returns a Kubernetes client that works both in-cluster and from kubeconfig
func NewClient() (*KubeClient, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fallback to local kubeconfig
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create k8s config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	return &KubeClient{Clientset: clientset}, nil
}

// GetKubeContext gathers basic cluster info: nodes and pods
func (kc *KubeClient) GetKubeContext() (map[string]interface{}, error) {
	ctx := context.TODO()

	nodes, err := kc.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := kc.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeNames := []string{}
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	podSummaries := []map[string]string{}
	for _, pod := range pods.Items {
		podSummaries = append(podSummaries, map[string]string{
			"name":      pod.Name,
			"namespace": pod.Namespace,
			"status":    string(pod.Status.Phase),
		})
	}

	result := map[string]interface{}{
		"node_count": len(nodeNames),
		"nodes":      nodeNames,
		"pod_count":  len(podSummaries),
		"pods":       podSummaries,
	}

	return result, nil
}
