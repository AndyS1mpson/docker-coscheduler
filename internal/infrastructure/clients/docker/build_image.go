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
		Tags:       []string{targetPath},
		Remove:     false,
	}

	res, err := c.externalClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer res.Body.Close()

	imageInspect, err := c.externalClient.InspectImage(targetPath)
	if err != nil {
		return "", fmt.Errorf("inspect image: %w", err)
	}

	return imageInspect.ID, nil
}
