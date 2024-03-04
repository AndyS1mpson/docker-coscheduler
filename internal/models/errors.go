package models

import "errors"

// DockerfileNotExistError ошибка существования Dockerfile'а
var DockerfileNotExistError error = errors.New("dockerfile does not exist")

// EmptyArchiveError ошибка пустого архива
var EmptyArchiveError error = errors.New("empty archive")
