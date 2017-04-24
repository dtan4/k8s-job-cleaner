# k8s-job-cleaner

Clean up succeeded Kubernetes Jobs

```bash
$ k8s-job-cleaner --label-group job --max-count 10
```

## Requirements

Kubernetes 1.3 or above

## Installation

TBD

## Usage

### In Kubernetes cluster

Just add `--in-cluster` flag.

```bash
$ k8s-job-cleaner --label-group job --max-count 10 --in-cluster
```

(TBD: CronJob manifest sample)

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
