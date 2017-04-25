package kubernetes

import (
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

// Jobs represents job list
type Jobs []batchv1.Job

// Len return the length of job list
func (j Jobs) Len() int {
	return len(j)
}

// Less returns whether the former item is less than the latter item or not
func (j Jobs) Less(m, n int) bool {
	return j[m].Status.CompletionTime.Before(*j[n].Status.CompletionTime)
}

// Swap swaps two items
func (j Jobs) Swap(m, n int) {
	j[m], j[n] = j[n], j[m]
}

// Jobs represents job list
type Pods []v1.Pod

// Sorting Pods is not necessary
