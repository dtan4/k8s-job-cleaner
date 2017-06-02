package kubernetes

import (
	"strings"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

func TestDeleteJob(t *testing.T) {
	obj := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "job",
			Namespace: "namespace",
		},
	}
	clientset := fake.NewSimpleClientset(obj)
	client := &Client{
		clientset: clientset,
	}

	testcases := []struct {
		job        batchv1.Job
		errMessage string
	}{
		{
			job: batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job",
					Namespace: "namespace",
				},
			},
			errMessage: "",
		},
		{
			job: batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "foobar",
					Namespace: "namespace",
				},
			},
			errMessage: "failed to delete Job",
		},
	}

	for _, testcase := range testcases {
		err := client.DeleteJob(testcase.job)

		if testcase.errMessage == "" {
			if err != nil {
				t.Errorf("error should not be raised: %s", err)
			}
		} else {
			if !strings.Contains(err.Error(), testcase.errMessage) {
				t.Errorf("error message does not contain %q; got: %q", testcase.errMessage, err.Error())
			}
		}
	}
}
