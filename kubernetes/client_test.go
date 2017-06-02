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

	if err := client.DeleteJob(batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "job",
			Namespace: "namespace",
		},
	}); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	err := client.DeleteJob(batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "foobar",
			Namespace: "namespace",
		},
	})

	if err == nil {
		t.Errorf("error should be raised")
	}

	expected := "failed to delete Job"

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error message does not contain %q; got: %q", expected, err.Error())
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

	if err := client.DeletePod(v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "pod",
			Namespace: "namespace",
		},
	}); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	err := client.DeletePod(v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "foobar",
			Namespace: "namespace",
		},
	})

	if err == nil {
		t.Errorf("error should be raised")
	}

	expected := "failed to delete Pod"

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error message does not contain %q; got: %q", expected, err.Error())
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
		namespace string
		expected  int
	}{
		{
			namespace: "namespace",
			expected:  2,
		},
		{
			namespace: "foobar",
			expected:  0,
		},
	}

	for _, testcase := range testcases {
		jobs, err := client.ListJobs(testcase.namespace)

		if err != nil {
			t.Errorf("error should not be raised: %s", err)
		}

		if len(jobs.Items) != testcase.expected {
			t.Errorf("the number of items does not match; expected: %d, got: %v", testcase.expected, jobs)
		}
	}
}

func TestListPods(t *testing.T) {
	obj := &v1.PodList{
		Items: []v1.Pod{
			v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod1",
					Namespace: "namespace",
				},
			},
			v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod2",
					Namespace: "namespace",
				},
			},
			v1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod3",
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
		namespace string
		expected  int
	}{
		{
			namespace: "namespace",
			expected:  2,
		},
		{
			namespace: "foobar",
			expected:  0,
		},
	}

	for _, testcase := range testcases {
		jobs, err := client.ListPods(testcase.namespace)

		if err != nil {
			t.Errorf("error should not be raised: %s", err)
		}

		if len(jobs.Items) != testcase.expected {
			t.Errorf("the number of items does not match; expected: %d, got: %v", testcase.expected, jobs)
		}
	}
}
