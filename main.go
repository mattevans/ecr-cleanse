package main

import (
	"flag"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/mattevans/ecr-cleaner-master/services"
)

var log = logrus.New()
var logGroup = []string{}

func main() {
	// Assign our flags.
	var (
		region = flag.String("aws-region", "us-west-2", "AWS region")
		dry    = flag.Bool("dry-run", false, "Executes without deleting any images")
		err    error
	)

	// Parse our flags.
	flag.Parse()

	// Intialize ECS and ECR clients
	ecsc := services.NewECSClient(*region)
	ecrc := services.NewECRClient(*region)

	// Find our repositories.
	repositories, err := ecrc.ListRepositories()
	if err != nil {
		log.WithError(err).Error("Error listing ECR repositories")
	}

	// Find image identities currenlty in use on active tasks.
	runningImages, err := ecsc.FindActiveImages()
	if err != nil {
		log.WithError(err).Error("Error finding `running` tasks within clusters")
	}

	// Log some info for the user.
	logGroup = append(logGroup, []string{
		fmt.Sprintf("Dry Run: %v", *dry),
		fmt.Sprintf("AWS Region: %v", *region),
		fmt.Sprintf("Repositories Found: %v", len(repositories)),
		fmt.Sprintf("Active Images Found: %v", len(runningImages)),
	}...)
	outputLog()

	// Range repo's, find stale images and remove them.
	for _, repository := range repositories {
		// Get our images within given repository.
		images, err := ecrc.ListImages(repository)
		if err != nil {
			log.Errorf("Error retrieving images for `%v` repository (%v)", repository, err)
		}
		logGroup = append(logGroup, fmt.Sprintf("Repository: %v", repository))

		// Build a slice of stale images.
		stale := []*ecr.ImageIdentifier{}
		for _, image := range images {
			match := false
			for _, running := range runningImages {
				// Handle case where we have an image without a tag.
				if image.ImageTag == nil {
					break
				}
				// Tag matches an image running a container. Ignore!
				if running == *image.ImageTag {
					match = true
					break
				}
			}
			// No match, image is stale.
			if match == false {
				stale = append(stale, image)
			}
		}

		// Purge our stale images if not a dry-run.
		if *dry {
			logGroup = append(logGroup, fmt.Sprintf("[DRY RUN] `%v` images would be purged", len(stale)))
		} else {
			err = ecrc.PurgeImages(repository, stale)
			if err != nil {
				log.Errorf("Error purging stale images for repo %v: %v", repository, err)
			}
			logGroup = append(logGroup, fmt.Sprintf("[PURGED] `%v` images in repository `%v`", len(stale), repository))
		}
		outputLog()
	}
}

func outputLog() {
	for _, value := range logGroup {
		log.Info(value)
	}
	log.Info("----------------------------------------------------------------")
	logGroup = []string{}
}
