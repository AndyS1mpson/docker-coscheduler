package image_hub

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnpackTarArchive(archiveFile *os.File, targetDir string) error {
	tarReader := tar.NewReader(archiveFile)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			fmt.Printf("Unsupported type: %v in %s\n", header.Typeflag, targetPath)
		}
	}

	return nil
}
