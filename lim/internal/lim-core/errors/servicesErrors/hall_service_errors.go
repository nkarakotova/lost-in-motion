package servicesErrors

import "errors"

var (
	HallAlreadyExists = errors.New("Service error! Зал уже существует в базе!")
	HallDoesNotExists = errors.New("Service error! Такого зала не существует!")
)
