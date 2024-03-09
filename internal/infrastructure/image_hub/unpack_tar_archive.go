package image_hub

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// UnpackTarArchive разархивирует архив и возвращает путь до папки с ним
func (h *Hub) UnpackTarArchive(archiveFile models.ImageArchive, dirName string) (string, error) {
	if archiveFile.File == nil {
		return "", models.ErrEmptyArchive
	}

	tarReader := tar.NewReader(archiveFile.File)

	targetPath := filepath.Join(h.workerImageHubDir, dirName)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}

		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return "", err
			}
		default:
			fmt.Printf("Unsupported type: %v in %s\n", header.Typeflag, targetPath)
		}
	}

	return targetPath, nil
}
