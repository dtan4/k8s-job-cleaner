package main

import (
	"fmt"
	"os"
	"sort"

	k8s "github.com/dtan4/k8s-job-cleaner/kubernetes"
	flag "github.com/spf13/pflag"
)

const (
	defaultMaxCount = 10
	jobNameLabel    = "job-name"
)

func main() {
	var (
		context    string
		dryRun     bool
		inCluster  bool
		kubeconfig string
		labelGroup string
		maxCount   int64
		namespace  string
		version    bool
	)

	flags := flag.NewFlagSet("k8stail", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.StringVar(&context, "context", "", "Kubernetes context")
	flags.BoolVar(&dryRun, "dry-run", false, "Dry run")
	flags.BoolVar(&inCluster, "in-cluster", false, "Execute in Kubernetes cluster")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVar(&labelGroup, "label-group", "", "Label name for grouping Jobs")
	flags.Int64Var(&maxCount, "max-count", int64(defaultMaxCount), "Number of Jobs to remain")
	flags.StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	flags.BoolVarP(&version, "version", "v", false, "Print version")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		} else {
			kubeconfig = k8s.DefaultConfigFile()
		}
	}

	if labelGroup == "" {
		fmt.Fprintln(os.Stderr, "--label-group must be set")
		os.Exit(1)
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	var client *k8s.Client

	if inCluster {
		c, err := k8s.NewClientInCluster()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if namespace == "" {
			namespace = k8s.DefaultNamespace()
		}

		client = c
	} else {
		c, err := k8s.NewClient(kubeconfig, context)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if namespace == "" {
			namespaceInConfig, err := c.NamespaceInConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if namespaceInConfig == "" {
				namespace = k8s.DefaultNamespace()
			} else {
				namespace = namespaceInConfig
			}
		}

		client = c
	}

	jobs, err := client.ListJobs(namespace)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jobGroup := map[string]k8s.Jobs{}

	for _, job := range jobs.Items {
		if job.Status.Succeeded == 0 {
			continue
		}

		label := job.Labels[labelGroup]
		if label == "" {
			continue
		}

		if jobGroup[label] == nil {
			jobGroup[label] = k8s.Jobs{}
		}

		jobGroup[label] = append(jobGroup[label], job)
	}

	pods, err := client.ListPods(namespace)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	podGroup := map[string]k8s.Pods{}

	for _, pod := range pods.Items {
		if !k8s.IsPodFinished(pod) {
			continue
		}

		label := pod.Labels[jobNameLabel]
		if label == "" {
			continue
		}

		if podGroup[label] == nil {
			podGroup[label] = k8s.Pods{}
		}

		podGroup[label] = append(podGroup[label], pod)
	}

	for _, jobs := range jobGroup {
		i := int64(0)
		sort.Sort(sort.Reverse(jobs))

		for _, job := range jobs {
			if i < maxCount {
				i++
				continue
			}

			if dryRun {
				fmt.Printf("Deleting Job %s... [dry-run]\n", job.Name)
			} else {
				fmt.Printf("Deleting Job %s...\n", job.Name)
				if err := client.DeleteJob(job); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			for _, pod := range podGroup[job.Name] {
				if dryRun {
					fmt.Printf("  Deleting Pod %s... [dry-run]\n", pod.Name)
				} else {
					fmt.Printf("  Deleting Pod %s...\n", pod.Name)
					if err := client.DeletePod(pod); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
			}
		}
	}
}
