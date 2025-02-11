package k8s

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Execute(language, code string, resultChan chan<- string, ctx context.Context) {
	log.Println("executed haha")
	Client := kubernetesClient
	if Client == nil {
		resultChan <- "Kubernetes client not initialized"
		return
	}

	// Determine filename and execution command based on language
	var filename, execCmd, imageName string
	switch language {
	case "java":
		filename = "Main.java"
		execCmd = fmt.Sprintf("echo '%s' > /%s && javac %s && java Main", code, filename, filename)
		imageName = "java-image:1.0.0"

	case "cpp":
		filename = "main.cpp"
		execCmd = fmt.Sprintf("echo '%s' > /%s && g++ %s -o main && ./main", code, filename, filename)
		imageName = "cpp-image"
	case "python":
		filename = "index.py"
		execCmd = fmt.Sprintf("echo '%s' > /%s && python %s", code, filename, filename)
		imageName = "python-image"
	case "golang":
		filename = "main.go"
		execCmd = fmt.Sprintf("echo '%s' > /%s && go run %s", code, filename, filename)
		imageName = "golang-image"
	case "c":
		filename = "main.c"
		execCmd = fmt.Sprintf("echo '%s' > /%s && gcc -o output_file %s && ./output_file", code, filename, filename)
		imageName = "c-image"
	case "nodejs":
		filename = "index.js"
		execCmd = fmt.Sprintf("echo '%s' > /%s && node %s", code, filename, filename)
		imageName = "node-image"
	default:
		resultChan <- fmt.Sprintf("unsupported language: %s", language)
		return
	}

	// Define the Kubernetes job
	jobName := fmt.Sprintf("code-executor-%d", time.Now().UnixNano())
	jobSpecs := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "executor",
							Image:           imageName,
							ImagePullPolicy: corev1.PullNever,
							Command: []string{
								"sh",
								"-c",
								execCmd,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	job, err := Client.BatchV1().Jobs("default").Create(ctx, jobSpecs, metav1.CreateOptions{})
	if err != nil {
		resultChan <- fmt.Sprintf("failed to create job: %v", err)
		return
	}
	log.Printf("Job %s created successfully", job.Name)
	// Wait for the job to complete
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			resultChan <- "execution cancelled"
			return
		case <-ticker.C:
			job, err := Client.BatchV1().Jobs("default").Get(ctx, job.Name, metav1.GetOptions{})
			if err != nil {
				resultChan <- fmt.Sprintf("failed to get job status: %v", err)
				return
			}

			if job.Status.Succeeded > 0 {
				log.Printf("Job %s completed successfully", job.Name)
				goto Logs
			}
		}
	}

Logs:
	// Get logs from the pod
	pods, err := Client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
	})
	if err != nil {
		resultChan <- fmt.Sprintf("failed to list pods: %v", err)
		return
	}

	if len(pods.Items) == 0 {
		resultChan <- "no pods found for the job"
		return
	}

	podName := pods.Items[0].Name
	logOptions := &corev1.PodLogOptions{Container: "executor"}
	logStream, err := Client.CoreV1().Pods("default").GetLogs(podName, logOptions).Stream(ctx)
	if err != nil {
		resultChan <- fmt.Sprintf("failed to get logs: %v", err)
		return
	}
	defer logStream.Close()

	var logOutput strings.Builder
	scanner := bufio.NewScanner(logStream)
	for scanner.Scan() {
		logOutput.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		resultChan <- fmt.Sprintf("error reading logs: %v", err)
		return
	}

	resultChan <- logOutput.String()

	// Clean up the job after completion
	err = Client.BatchV1().Jobs("default").Delete(ctx, job.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("failed to delete job: %v", err)
	}
}
