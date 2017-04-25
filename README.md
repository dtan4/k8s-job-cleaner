# k8s-job-cleaner

[![Build Status](https://travis-ci.org/dtan4/k8s-job-cleaner.svg?branch=master)](https://travis-ci.org/dtan4/k8s-job-cleaner)
[![Docker Repository on Quay](https://quay.io/repository/dtan4/k8s-job-cleaner/status "Docker Repository on Quay")](https://quay.io/repository/dtan4/k8s-job-cleaner)

Clean up completed Kubernetes Jobs

## What is this?

This tool provides Job cleaning feature based on k8s 1.6's [Job History Limits](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#jobs-history-limits).

For example, following command deletes completed Jobs and attached Pods, but leaves the last 10 Jobs per `job` label.

```bash
$ k8s-job-cleaner --label-group job --max-count 10
```

## Requirements

Kubernetes 1.3 or above

## Installation

### From source

```bash
$ go get -d github.com/dtan4/k8s-job-cleaner
$ cd $GOPATH/src/github.com/dtan4/k8s-job-cleaner
$ make deps
$ make install
```

### Run in a Docker container

Docker image is available at [quay.io/dtan4/k8s-job-cleaner](https://quay.io/repository/dtan4/k8s-job-cleaner).

```bash
# -t is required to colorize logs
$ docker run \
    --rm \
    -t \
    -v $HOME/.kube/config:/.kube/config \
    quay.io/dtan4/k8s-job-cleaner:latest \
      --label-group job
```

## Usage

### In Kubernetes cluster

Just add `--in-cluster` flag.

```bash
$ k8s-job-cleaner --label-group job --max-count 10 --in-cluster
```

CronJob manifest sample:

```yaml
apiVersion: batch/v2alpha1
kind: CronJob
metadata:
  name: k8s-job-cleaner
  labels:
    job: k8s-job-cleaner
    role: job
spec:
  schedule: "0 * * * *"
  startingDeadlineSeconds: 30
  concurrencyPolicy: Allow
  suspend: false
  jobTemplate:
    metadata:
      name: k8s-job-cleaner
      labels:
        job: k8s-job-cleaner
        role: job
    spec:
      template:
        metadata:
          name: k8s-job-cleaner
          labels:
            job: k8s-job-cleaner
            role: job
        spec:
          containers:
          - name: k8s-job-cleaner
            image: quay.io/dtan4/k8s-job-cleaner:latest
            imagePullPolicy: Always
            command:
              - "/k8s-job-cleaner"
              - "--in-cluster"
              - "--label-group"
              - "job"
              - "--max-count"
              - "20"
          restartPolicy: Never
```

### Local machine

`k8s-job-cleaner` uses `~/.kube/config` as default. You can specify another path by `KUBECONFIG` environment variable or `--kubeconfig` option. `--kubeconfig` option always overrides `KUBECONFIG` environment variable.

```bash
$ KUBECONFIG=/path/to/kubeconfig k8s-job-cleaner
# or
$ k8s-job-cleaner --kubeconfig=/path/to/kubeconfig
```

### Options

|Option|Description|Required|Default|
|---------|-----------|-------|-------|
|`--context=CONTEXT`|Kubernetes context|||
|`--dry-run`|Dry run||`false`|
|`--in-cluster`|Execute in Kubernetes cluster||`false`|
|`--kubeconfig=KUBECONFIG`|Path of kubeconfig||`~/.kube/config`|
|`--label-group=LABELS`|Label name for groupiung Jobs|Required||
|`--max-count=MAXCOUNT`|Number of Jobs to remain||`10`|
|`--namespace=NAMESPACE`|Kubernetes namespace||All namespaces|
|`-h`, `-help`|Print command line usage|||
|`-v`, `-version`|Print version|||

## Development

Clone this repository and build using `make`.

```bash
$ go get -d github.com/dtan4/k8s-job-cleaner
$ cd $GOPATH/src/github.com/dtan4/k8s-job-cleaner
$ make
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
