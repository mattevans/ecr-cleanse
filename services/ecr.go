package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECRClient holds a connection to AWS ECR.
type ECRClient struct {
	client *ecr.ECR
	region string
}

// NewECRClient initializes an ECRClient.
func NewECRClient(region string) *ECRClient {
	return &ECRClient{
		client: ecr.New(session.New(), aws.NewConfig().WithRegion(region)),
		region: region,
	}
}

// ListRepositories will return all ECR repositories as []string.
func (c *ECRClient) ListRepositories() ([]string, error) {
	allRepositories, err := c.client.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		return []string{}, err
	}

	repositories := make([]string, 0, len(allRepositories.Repositories))
	for _, repo := range allRepositories.Repositories {
		repositories = append(repositories, *repo.RepositoryName)
	}
	return repositories, nil
}

// ListImages will return all image identifiers for a given repository.
func (c *ECRClient) ListImages(repository string) ([]*ecr.ImageIdentifier, error) {
	var token *string
	var imageIDs []*ecr.ImageIdentifier

	for {
		resp, err := c.client.ListImages(&ecr.ListImagesInput{
			RepositoryName: aws.String(repository),
			NextToken:      token,
		})
		if err != nil {
			return nil, err
		}

		imageIDs = append(imageIDs, resp.ImageIds...)
		if token = resp.NextToken; token == nil {
			break
		}
	}
	return imageIDs, nil
}

// PurgeImages will batch delete images by image identitfier in sets of 100.
func (c *ECRClient) PurgeImages(repository string, images []*ecr.ImageIdentifier) error {
	// No images found, back-out.
	if len(images) <= 0 {
		return nil
	}
	// Purge the images in batches of 100.
	i := 0
	for i = 0; i < int(len(images)/100); i++ {
		err := c.BatchPurge(repository, images[i*100:(i+1)*100])
		if err != nil {
			return fmt.Errorf("Failed purging images in repository `%v` (%v)", repository, err)
		}
	}
	err := c.BatchPurge(repository, images[i*100:])
	if err != nil {
		return fmt.Errorf("Failed purging images in repository `%v` (%v)", repository, err)
	}
	return err
}

// BatchPurge will batch delete images by the image identifiers.
func (c *ECRClient) BatchPurge(repository string, images []*ecr.ImageIdentifier) error {
	_, err := c.client.BatchDeleteImage(&ecr.BatchDeleteImageInput{
		RepositoryName: aws.String(repository),
		ImageIds:       images,
	})
	if err != nil {
		return err
	}
	return nil
}
