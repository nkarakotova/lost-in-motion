package servicesErrors

import "errors"

var (
	NoAvailablePlacesNum      = errors.New("Service error! Не осталось свободных мест на тренировке!")
	IncorrectTrainingTime     = errors.New("Service error! Тренировка начиниется в недопустимое время!")
	BusyDateTime              = errors.New("Service error! В данное время занят зал или тренер!")
	StartDateAfterEndDate     = errors.New("Service error! Начальная дата позже конечной!")
	TrainingDoesNotExists     = errors.New("Service error! Такой тренировки не существует!")
)
