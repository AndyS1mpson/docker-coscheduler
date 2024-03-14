package image_hub

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var (
	dockerfileName = "Dockerfile"

	tarExtension = ".tar"
)

// ArchiveImageToTar создает tar архив образа Docker
func (h *Hub) ArchiveImageToTar(imageDir string, tarName string) (*models.ImageArchive, error) {
	if !isDockerfileExist(filepath.Join(h.schedulerTaskDir, imageDir)) {
		return nil, models.ErrDockerfileNotExist
	}

	archive, err := h.createTarArchive(imageDir, tarName)
	if err != nil {
		return nil, fmt.Errorf("create tar archive: %w", err)
	}

	return archive, nil
}

// createTarArchive создает tar архив
func (h *Hub) createTarArchive(imageDir string, tarName string) (*models.ImageArchive, error) {
	folderPath := filepath.Join(h.schedulerTaskDir, imageDir)

	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)
	defer tw.Close()

	// Функция для обхода файлов и папок внутри указанной директории
	walker := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	}

	if err := filepath.Walk(folderPath, walker); err != nil {
		return nil, err
	}

	return &models.ImageArchive{File: buf.Bytes()}, nil
}

// isDockerfileExist проверяет, содержит ли указанный каталог файл Dockerfile
func isDockerfileExist(imageDir string) bool {
	_, err := os.Stat(filepath.Join(imageDir, dockerfileName))

	return err == nil
}
