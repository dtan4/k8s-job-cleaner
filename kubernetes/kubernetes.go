package kubernetes

import (
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// DefaultConfigFile returns the default kubeconfig file path
func DefaultConfigFile() string {
	return clientcmd.RecommendedHomeFile
}

// DefaultNamespace returns the default namespace
func DefaultNamespace() string {
	return v1.NamespaceAll
}

// IsJobFinished returns whether the given Job has finished or not
func IsJobFinished(job batchv1.Job) bool {
	return job.Status.Succeeded > 0
}

// IsPodFinished returns whether the given Pod has finished or not
func IsPodFinished(pod v1.Pod) bool {
	return pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed
}
