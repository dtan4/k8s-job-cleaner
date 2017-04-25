package kubernetes

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client represents the wrapper of Kubernetes API client
type Client struct {
	clientConfig clientcmd.ClientConfig
	clientset    *kubernetes.Clientset
}

// DefaultConfigFile returns the default kubeconfig file path
func DefaultConfigFile() string {
	return clientcmd.RecommendedHomeFile
}

// DefaultNamespace returns the default namespace
func DefaultNamespace() string {
	return v1.NamespaceAll
}

// IsPodFinished returns whether the given Pod has finished or has not
func IsPodFinished(pod v1.Pod) bool {
	return pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed
}

// NewClient creates Client object using local kubecfg
func NewClient(kubeconfig, context string) (*Client, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: context})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "falied to load local kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load clientset")
	}

	return &Client{
		clientConfig: clientConfig,
		clientset:    clientset,
	}, nil
}

// NewClientInCluster creates Client object in Kubernetes cluster
func NewClientInCluster() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubeconfig in cluster")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "falied to load clientset")
	}

	return &Client{
		clientset: clientset,
	}, nil
}

// DeleteJob deletes the given Job
func (c *Client) DeleteJob(job batchv1.Job) error {
	if err := c.clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name, &v1.DeleteOptions{}); err != nil {
		return errors.Wrap(err, "failed to delete Job")
	}

	return nil
}

// DeletePod deletes the given Pod
func (c *Client) DeletePod(pod v1.Pod) error {
	if err := c.clientset.Core().Pods(pod.Namespace).Delete(pod.Name, &v1.DeleteOptions{}); err != nil {
		return errors.Wrap(err, "failed to delete Pod")
	}

	return nil
}

// ListJobs returns the list of Jobs
func (c *Client) ListJobs(namespace string) (*batchv1.JobList, error) {
	jobs, err := c.clientset.BatchV1().Jobs(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve Jobs")
	}

	return jobs, nil
}

// ListPods returns the list of Pods
func (c *Client) ListPods(namespace string) (*v1.PodList, error) {
	pods, err := c.clientset.Core().Pods(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve Pods")
	}

	return pods, nil
}

// NamespaceInConfig returns namespace set in kubeconfig
func (c *Client) NamespaceInConfig() (string, error) {
	if c.clientConfig == nil {
		return "", errors.New("clientConfig is not set")
	}

	rawConfig, err := c.clientConfig.RawConfig()
	if err != nil {
		return "", errors.Wrap(err, "failed to load rawConfig")
	}

	return rawConfig.Contexts[rawConfig.CurrentContext].Namespace, nil
}
