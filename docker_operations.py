package main

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
)

// DockerOperations represents Docker operations
type DockerOperations struct {
    client *client.Client
}

// NewDockerOperations creates a new instance of DockerOperations
func NewDockerOperations() (*DockerOperations, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return nil, err
    }
    return &DockerOperations{client: cli}, nil
}

// ExecuteScript executes a script in a Docker container
func (d *DockerOperations) ExecuteScript(image, script string) (string, error) {
    ctx := context.Background()
    resp, err := d.client.ContainerCreate(ctx, &container.Config{
        Image: image,
        Cmd:   []string{"bash", "-c", script},
    }, nil, nil, nil, "")
    if err != nil {
        return "", fmt.Errorf("failed to create container: %v", err)
    }
    defer d.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})

    if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        return "", fmt.Errorf("failed to start container: %v", err)
    }

    statusCh, errCh := d.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            return "", fmt.Errorf("failed to wait for container: %v", err)
        }
    case <-statusCh:
    }

    out, err := d.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
    if err != nil {
        return "", fmt.Errorf("failed to get container logs: %v", err)
    }
    defer out.Close()

    output, err := io.ReadAll(out)
    if err != nil {
        return "", fmt.Errorf("failed to read container logs: %v", err)
    }

    return string(output), nil
}