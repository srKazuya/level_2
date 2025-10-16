package cmdErrs

import (
	"errors"
)

var (
	ErrExtCommand = errors.New("Внешняя ошибка: \n")
	ErrGetwd     = errors.New("Ошибка при получении директории:\n")
	ErrChangeDir = errors.New("Ошибка смены директории:\n")
)
