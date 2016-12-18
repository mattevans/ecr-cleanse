package services

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSClient holds a connection to AWS ECS.
type ECSClient struct {
	client *ecs.ECS
	region string
}

// NewECSClient initializes an ECSClient.
func NewECSClient(region string) *ECSClient {
	return &ECSClient{
		client: ecs.New(session.New(), aws.NewConfig().WithRegion(region)),
		region: region,
	}
}

// FindActiveImages will initiate a connection with ECS, parsing out image names
// from all running tasks/clusters, returning them as a []string.
func (c *ECSClient) FindActiveImages() ([]string, error) {
	// List ECS clusters.
	var token *string
	clusters, err := c.client.ListClusters(&ecs.ListClustersInput{
		NextToken: token,
	})
	if err != nil {
		return nil, err
	}

	// Slice to store our running images name(s).
	runningImages := []string{}

	// Range cluster ARNs and parse out the image ID from the running task.
	clusterARNs := make([]string, 0, len(clusters.ClusterArns))
	for _, arn := range clusters.ClusterArns {
		clusterARNs = append(clusterARNs, *arn)

		// List all 'running' tasks from cluster.
		running := ecs.DesiredStatusRunning
		runningTasks, err := c.client.ListTasks(&ecs.ListTasksInput{
			Cluster:       arn,
			DesiredStatus: &running,
		})
		if err != nil {
			return nil, err
		}

		// Describe cluster tasks.
		tasks, err := c.client.DescribeTasks(&ecs.DescribeTasksInput{
			Tasks:   runningTasks.TaskArns,
			Cluster: arn,
		})
		if err != nil {
			return nil, err
		}

		// Range the tasks, parsing/storing the running ECI name.
		for _, task := range tasks.Tasks {
			tasks, err := c.client.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
				TaskDefinition: task.TaskDefinitionArn,
			})
			if err != nil {
				return nil, err
			}
			for _, container := range tasks.TaskDefinition.ContainerDefinitions {
				runningImages = append(runningImages, strings.Split(*container.Image, ":")[1])
			}
		}
	}
	return runningImages, nil
}
