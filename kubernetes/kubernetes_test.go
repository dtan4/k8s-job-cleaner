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
			t.Errorf("result is wrong: expected: %b", testcase.expected)
		}
	}
}
