package models

import "errors"

// DockerfileNotExistError ошибка существования Dockerfile'а
var DockerfileNotExistError error = errors.New("dockerfile does not exist")
