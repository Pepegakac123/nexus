package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Pepegakac123/nexus/pkg/orchestrator"
)

func main() {
	ctx := context.Background()
	o, err := orchestrator.NewOrchestrator(ctx)
	if err != nil {
		fmt.Printf("Error creating orchestrator: %v\n", err)
		os.Exit(1)
	}
	// List containers
	result, err := o.ListContainers(ctx)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s  %-22s  %s\n", "ID", "STATUS", "IMAGE")
	for _, ctr := range result.Items {
		fmt.Printf("%s  %-22s  %s\n", ctr.ID, ctr.Status, ctr.Image)
	}
}
