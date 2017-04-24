package main

import (
	"fmt"
)

var (
	// Version represents the version number
	Version string
	// Revision represents the Git commit ID as of build
	Revision string
)

func printVersion() {
	fmt.Printf("k8s-job-cleaner version %s, %s\n", Version, Revision)
}
