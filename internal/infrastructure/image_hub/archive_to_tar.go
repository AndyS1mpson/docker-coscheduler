package image_hub

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var (
	dockerfileName = "Dockerfile"

	tarExtension   = ".tar"
	tarArchivesDir = "/tmp"
)

// ArchiveImageToTar создает tar архив образа Docker
func (h *Hub) ArchiveImageToTar(imageDir string, tarName string) (*models.ImageArchive, error) {
	if !isDockerfileExist(imageDir) {
		return nil, models.ErrDockerfileNotExist
	}

	tarFile, err := createTarArchive(imageDir, tarName)
	if err != nil {
		return nil, fmt.Errorf("create tar archive: %w", err)
	}

	return &models.ImageArchive{File: tarFile}, nil
}

// isDockerfileExist проверяет, содержит ли указанный каталог файл Dockerfile
func isDockerfileExist(imageDir string) bool {
	_, err := os.Stat(filepath.Join(imageDir, dockerfileName))

	return err == nil
}

// createTarArchive создает tar архив
func createTarArchive(imageDir string, tarName string) (*os.File, error) {
	tarPath := filepath.Join(tarArchivesDir, tarName+tarExtension)

	tarFile, err := os.Create(tarPath)
	if err != nil {
		return nil, err
	}

	tarWriter := tar.NewWriter(tarFile)

	err = filepath.Walk(imageDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(imageDir, path)
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = tarWriter.Close()
	if err != nil {
		return nil, err
	}

	return tarFile, nil
}
