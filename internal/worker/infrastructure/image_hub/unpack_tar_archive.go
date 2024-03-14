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

// UnpackTarArchive разархивирует архив и возвращает путь до папки с ним
func (h *Hub) UnpackTarArchive(archiveFile models.ImageArchive, dirName string) error {
	if archiveFile.File == nil {
		return models.ErrEmptyArchive
	}

	buf := bytes.NewBuffer(archiveFile.File)

	tarReader := tar.NewReader(buf)

	folderPath := filepath.Join(h.imageHubDir, dirName)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("create folder: %w", err)
	}

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF: // больше нет файлов
			return nil
		case err != nil: // ошибка чтения
			return err
		case header == nil: // скип файла если header nil
			continue
		}

		target := filepath.Join(folderPath, header.Name)

		// Распаковка и сохранение файлов из архива
		switch header.Typeflag {
		// Если это директория
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// Если это файл
		case tar.TypeReg:
			// Создание файла в файловой системе
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			defer outFile.Close()

			// Запись содержимого файла из архива в файловую систему
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
		}

	}
}
