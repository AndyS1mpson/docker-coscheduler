package models

import "os"

// ImageArchive представляет архив образа Docker
type ImageArchive struct {
	File *os.File
}
