package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/Pepegakac123/nexus/pkg/orchestrator"
	"github.com/moby/moby/client"
)

func main() {
	workers := runtime.NumCPU()
	ctx := context.Background()
	o, err := orchestrator.NewOrchestrator(ctx)
	if err != nil {
		fmt.Printf("Error creating orchestrator: %v\n", err)
		os.Exit(1)
	}
	listResult, err := o.ListContainers(ctx)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		os.Exit(1)
	}
	wg := sync.WaitGroup{}
	jobs := make(chan string, workers)
	results := make(chan client.ContainerInspectResult, workers)
	for range workers {
		wg.Go(func() {
			o.Worker(ctx, jobs, results)
		})

	}
	go func() { wg.Wait(); close(results) }()
	go func() {
		for _, ctr := range listResult.Items {
			jobs <- ctr.ID
		}
		close(jobs)
	}()
	fmt.Printf("%-12s %-20s %s\n", "ID", "IMAGE", "STATUS")
	for res := range results {
		if res.Container.ID == "" {
			continue
		}
		shortID := res.Container.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}
		fmt.Printf("%-12s %-20s %s\n", shortID, res.Container.Image, res.Container.State.Status)
	}

}
