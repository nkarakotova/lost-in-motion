package servicesErrors

import "errors"

var (
	ClientDoesNotExists = errors.New("Service error! Такого клиента не существует!")
	ClientAlreadyExists = errors.New("Service error! Клиент уже существует в базе!")
	InvalidPassword = errors.New("Service error! Неверный пароль!")
	AssignmentOnThisTimeAlreadyExists = errors.New("Service error! Клиент уже записан на тренировку в это время!")
	ClientTelephoneIncorrect = errors.New("Service error! Некорректный телефон клиента!")
	ClientMailIncorrect = errors.New("Service error! Некорректная почта клиента!")
	AssignmentDoesNotExists = errors.New("Service error! Запись на тренировку не существует!")
)
