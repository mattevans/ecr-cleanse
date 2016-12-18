# ecr-cleanse

A golang implementation for removing old images from ECR repositories.

This script will inspect all running tasks, across all ECS clusters, ensuring only images that are not in-use (`RUNNING`) are removed.§

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
