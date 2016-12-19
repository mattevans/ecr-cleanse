# ecr-cleanse

A golang implementation for removing old images from ECR repositories.

This script will inspect all services/tasks, across all ECS clusters, removing only those images that are not in-use.

Usage
-----------------

```go
go run main.go
```

Flags
-----------------

- `-aws-region`: Pass your AWS region.
- `-dry-run`: Execute the script without purging any images.

```go
go run main.go -aws-region us-west-2 -dry-run
```

Example Output
-----------------

```
INFO[0005] Dry Run: true
INFO[0005] AWS Region: us-west-2
INFO[0005] Repositories Found: 3
INFO[0005] Active Images Found: 8
INFO[0005] ----------------------------------------------------------------
INFO[0005] Repository: my.production.repository
INFO[0005] [DRY RUN] `2` images would be purged
INFO[0005] ----------------------------------------------------------------
INFO[0005] Repository: my.staging.repository
INFO[0005] [DRY RUN] `2` images would be purged
INFO[0005] ----------------------------------------------------------------
INFO[0006] Repository: my.test.repository
INFO[0006] [DRY RUN] `1` images would be purged
INFO[0006] ----------------------------------------------------------------
```

Contributing
-----------------
If you've found a bug or would like to contribute, please create an issue here on GitHub, or better yet fork the project and submit a pull request!
