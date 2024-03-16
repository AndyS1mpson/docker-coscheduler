package docker

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

var dockerfileName = "Dockerfile"

// BuildImage сборка докер образа и возвращение идентификатора образа
func (c *Client) BuildImage(ctx context.Context, dirName string) (string, error) {
	targetPath := filepath.Join(c.workerImageHubDir, dirName)

	tar, err := archive.TarWithOptions(targetPath, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("unpack tar archive: %w", err)
	}

	opts := types.ImageBuildOptions{
		Dockerfile: dockerfileName,
		Tags:       []string{dirName},
		Remove:     true,
	}

	res, err := c.externalClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer res.Body.Close()

	images, err := c.externalClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return "", fmt.Errorf("inspect image: %w", err)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == fmt.Sprintf("%s:latest", dirName) {
				return image.ID[7:19], nil
			}
		}
	}

	return "", fmt.Errorf("image with tag %s not found", dirName)
}
