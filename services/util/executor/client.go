package executor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)



// RemoveContainerByID removes a container by its ID.
func RemoveContainerByID(client *client.Client, containerID string) error {
	// Forcefully remove the container
	if err := client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}
	log.Println(" removing ########### ", containerID)
	return nil
}

// RemoveAllContainers removes all containers.
func RemoveAllContainers(client *client.Client) error {
	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}
	for _, c := range containers {
		if err := RemoveContainerByID(client, c.ID); err != nil {
			return err
		}
	}
	return nil
}

func Execute(language string, code string, resultChan chan<- string, ctx context.Context) {
	client, err := GetClient()
	if err != nil {
		resultChan <- fmt.Sprintf("failed to get Docker client: %v", err)
		return
	}

	// Determine filename and execution command based on language
	var filename string
	var execCmd string
	var imageName string
	switch language {
	case "java":
		filename = "Main.java"
		execCmd = "javac Main.java && java Main"
		imageName = "java-image"
	case "cpp":
		filename = "main.cpp"
		execCmd = "g++ main.cpp -o main && ./main"
		imageName = "cpp-image"
	case "python":
		filename = "index.py"
		execCmd = "python index.py"
		imageName = "python-image"
	case "golang":
		filename = "main.go"
		execCmd = "go run main.go"
		imageName = "golang-image"
	case "c":
		filename = "main.c"
		execCmd = "gcc -o output_file main.c && ./output_file"
		imageName = "c-image"
	case "nodejs":
		filename = "index.js"
		execCmd = "node index.js"
		imageName = "node-image"
	default:
		resultChan <- fmt.Sprintf("unsupported language: %s", language)
		return
	}

	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo '%s' > /%s && cd / && %s", code, filename, execCmd)},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		resultChan <- fmt.Sprintf("failed to create container: %v", err)
		return
	}
	containerId := resp.ID

	// Cleanup the container when the context is done
	go func() {
		<-ctx.Done()
		RemoveContainerByID(client, containerId)
	}()

	// Start the container
	if err := client.ContainerStart(context.Background(), containerId, container.StartOptions{}); err != nil {
		resultChan <- fmt.Sprintf("failed to start container: %v", err)
		return
	}

	log.Println("Container started successfully", containerId)

	// Stream logs in a separate goroutine
	var logOutput strings.Builder
	go func() {
		out, err := client.ContainerLogs(context.Background(), containerId, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			resultChan <- fmt.Sprintf("failed to get container logs: %v", err)
			return
		}
		defer out.Close()

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			logOutput.WriteString(scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			resultChan <- fmt.Sprintf("error reading logs: %v", err)
		}
	}()

	// Wait for the container to finish executing
	statusCh, errCh := client.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			resultChan <- fmt.Sprintf("error waiting for container: %v", err)
			return
		}
	case <-statusCh:
		// Container finished, send the accumulated logs
		resultChan <- logOutput.String()
	}
}

// remove container periodically !!
func StartCleanupJob(client *client.Client) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
			if err != nil {
				log.Printf("failed to list containers: %v\n", err)
				continue
			}

			for _, c := range containers {
				// Check if container is older than 15 minutes
				if time.Since(time.Unix(c.Created, 0)) > 15*time.Minute {
					err := RemoveContainerByID(client, c.ID)
					if err != nil {
						log.Printf("failed to remove container %s: %v\n", c.ID, err)
					}
				}
			}
		}
	}()
}
