package main

import (
	"fmt"
	"os"
	"sort"

	k8s "github.com/dtan4/k8s-job-cleaner/kubernetes"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultMaxCount  = 10
	defaultNamespace = v1.NamespaceAll
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
			kubeconfig = clientcmd.RecommendedHomeFile
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
			namespace = defaultNamespace
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
				namespace = defaultNamespace
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

	jobGroup := map[string]Jobs{}

	for _, job := range jobs.Items {
		if job.Status.Succeeded == 0 {
			continue
		}

		label := job.Labels[labelGroup]
		if label == "" {
			continue
		}

		if jobGroup[label] == nil {
			jobGroup[label] = Jobs{}
		}

		jobGroup[label] = append(jobGroup[label], job)
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
		}
	}
}
