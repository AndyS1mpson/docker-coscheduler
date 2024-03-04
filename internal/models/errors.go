package models

import "errors"

// ErrDockerfileNotExist ошибка существования Dockerfile'а
var ErrDockerfileNotExist error = errors.New("dockerfile does not exist")

// ErrEmptyArchive ошибка пустого архива
var ErrEmptyArchive error = errors.New("empty archive")
