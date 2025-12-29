package orchestrator

import (
	"context"

	"github.com/moby/moby/client"
)

type Orchestrator struct {
	client *client.Client
}

func NewOrchestrator(ctx context.Context) (*Orchestrator, error) {
	client, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Orchestrator{
		client: client,
	}, nil
}

func (o *Orchestrator) ListContainers(ctx context.Context) (client.ContainerListResult, error) {
	result, err := o.client.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return client.ContainerListResult{}, err
	}
	return result, nil
}

func (o *Orchestrator) InspectContainer(ctx context.Context, id string) (client.ContainerInspectResult, error) {
	result, err := o.client.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return client.ContainerInspectResult{}, err
	}
	return result, nil
}

func (o *Orchestrator) Worker(ctx context.Context, jobs <-chan string, results chan<- client.ContainerInspectResult) {
	for id := range jobs {
		result, err := o.InspectContainer(ctx, id)
		if err != nil {
			results <- client.ContainerInspectResult{}
		} else {
			results <- result
		}
	}
}
