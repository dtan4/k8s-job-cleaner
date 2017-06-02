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
			if err == nil {
				t.Errorf("error should be raised")
			}

			if !strings.Contains(err.Error(), testcase.errMessage) {
				t.Errorf("error message does not contain %q; got: %q", testcase.errMessage, err.Error())
			}
		}
	}
}

func TestDeletePod(t *testing.T) {
	obj := &v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "pod",
			Namespace: "namespace",
		},
	}
	clientset := fake.NewSimpleClientset(obj)
	client := &Client{
		clientset: clientset,
	}

	testcases := []struct {
		pod        v1.Pod
		errMessage string
	}{
		{
			pod: v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod",
					Namespace: "namespace",
				},
			},
			errMessage: "",
		},
		{
			pod: v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "foobar",
					Namespace: "namespace",
				},
			},
			errMessage: "failed to delete Pod",
		},
	}

	for _, testcase := range testcases {
		err := client.DeletePod(testcase.pod)

		if testcase.errMessage == "" {
			if err != nil {
				t.Errorf("error should not be raised: %s", err)
			}
		} else {
			if err == nil {
				t.Errorf("error should be raised")
			}

			if !strings.Contains(err.Error(), testcase.errMessage) {
				t.Errorf("error message does not contain %q; got: %q", testcase.errMessage, err.Error())
			}
		}
	}
}

func TestListJobs(t *testing.T) {
	obj := &batchv1.JobList{
		Items: []batchv1.Job{
			batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job1",
					Namespace: "namespace",
				},
			},
			batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job2",
					Namespace: "namespace",
				},
			},
			batchv1.Job{
				ObjectMeta: v1.ObjectMeta{
					Name:      "job3",
					Namespace: "namespace-foobar",
				},
			},
		},
	}
	clientset := fake.NewSimpleClientset(obj)
	client := &Client{
		clientset: clientset,
	}

	testcases := []struct {
		namespace  string
		expected   int
		errMessage string
	}{
		{
			namespace:  "namespace",
			expected:   2,
			errMessage: "",
		},
		{
			namespace:  "foobar",
			expected:   0,
			errMessage: "",
		},
	}

	for _, testcase := range testcases {
		jobs, err := client.ListJobs(testcase.namespace)

		if testcase.errMessage == "" {
			if err != nil {
				t.Errorf("error should not be raised: %s", err)
			}

			if len(jobs.Items) != testcase.expected {
				t.Errorf("the number of items does not match; expected: %d, got: %v", testcase.expected, jobs)
			}
		} else {
			if err == nil {
				t.Errorf("error should be raised")
			}

			if !strings.Contains(err.Error(), testcase.errMessage) {
				t.Errorf("error message does not contain %q; got: %q", testcase.errMessage, err.Error())
			}
		}
	}
}
