package kube

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	Replicas  int32  `json:"replicas"`
	Port      int32  `json:"port"`
}

func ApplyDeployment(clientset *kubernetes.Clientset, spec DeploymentSpec) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: spec.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": spec.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": spec.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  spec.Name,
							Image: spec.Image,
							Ports: []corev1.ContainerPort{{ContainerPort: spec.Port}},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments(spec.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func DeleteDeployment(clientset *kubernetes.Clientset, name, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return clientset.AppsV1().Deployments(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
}
