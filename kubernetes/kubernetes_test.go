package kubernetes

import (
	"testing"

	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

func TestIsJobFinished(t *testing.T) {
	testcases := []struct {
		job      batchv1.Job
		expected bool
	}{
		{
			job: batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job-success",
					Namespace: "namespace",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			expected: true,
		},
		{
			job: batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job-failed",
					Namespace: "namespace",
				},
				Status: batchv1.JobStatus{
					Succeeded: 0,
				},
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		if IsJobFinished(testcase.job) != testcase.expected {
			t.Errorf("wrong result: expected: %b", testcase.expected)
		}
	}
}

func TestIsPodFinished(t *testing.T) {
	testcases := []struct {
		pod      v1.Pod
		expected bool
	}{
		{
			pod: v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod-success",
					Namespace: "namespace",
				},
				Status: v1.PodStatus{
					Phase: v1.PodSucceeded,
				},
			},
			expected: true,
		},
		{
			pod: v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod-failed",
					Namespace: "namespace",
				},
				Status: v1.PodStatus{
					Phase: v1.PodFailed,
				},
			},
			expected: true,
		},
		{
			pod: v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod-failed",
					Namespace: "namespace",
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
				},
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		if IsPodFinished(testcase.pod) != testcase.expected {
			t.Errorf("wrong result: expected: %b", testcase.expected)
		}
	}
}
