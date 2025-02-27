package servicesErrors

import "errors"

var (
	CoachAlreadyExists = errors.New("Service error! Тренер уже существует в базе!")
	CoachDoesNotExists = errors.New("Service error! Такого тренера не существует!")
)
