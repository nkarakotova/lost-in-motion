package techUI

import (
	"context"
	"fmt"

	"lim/internal/lim-console/registry"
	"lim/internal/lim-core/errors/menuErrors"
	"lim/internal/lim-core/models"
)

func loginClient(a *registry.AppServiceFields) (*models.Client, error) {
	var telephone string
	fmt.Printf("Введите телефон: ")
	_, err := fmt.Scanf("%s", &telephone)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var password string
	fmt.Printf("Введите пароль: ")
	_, err = fmt.Scanf("%s", &password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ctx := context.Background()

	client, err := a.ClientService.Login(ctx, telephone, password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}

func createClient(a *registry.AppServiceFields) (*models.Client, error) {
	var name string
	fmt.Printf("Введите имя: ")
	_, err := fmt.Scanf("%s", &name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var telephone string
	fmt.Printf("Введите телефон: ")
	_, err = fmt.Scanf("%s", &telephone)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var password string
	fmt.Printf("Введите пароль: ")
	_, err = fmt.Scanf("%s", &password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var mail string
	fmt.Printf("Введите почту: ")
	_, err = fmt.Scanf("%s", &mail)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &models.Client{Name: name,
		Telephone: telephone,
		Mail:      mail,
		Password:  password}

	ctx := context.Background()

	err = a.ClientService.Create(ctx, client)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}

func printInfoClient(client *models.Client) {
	fmt.Printf("\nИмя: %s\nТелефон: %s\nПочта: %s\n\n", client.Name, client.Telephone, client.Mail)
}

func CreateAssignment(a *registry.AppServiceFields, client *models.Client) {
	err := printTrainingsOnWeek(a)
	if err != nil {
		return
	}

	var tID uint64
	fmt.Printf("Выберете тренировку и введите её id: ")
	_, err = fmt.Scanf("%d", &tID)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	err = a.ClientService.CreateAssignment(ctx, client.ID, tID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Запись успешно создана!")
}

func printassignment(a *registry.AppServiceFields, client *models.Client) error {
	ctx := context.Background()

	trainings, err := a.TrainingService.GetAllByClient(ctx, client.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(trainings) == 0 {
		fmt.Println("Ещё нет записей на тренировки!")
		return menuErrors.ErrorMenu
	}

	fmt.Printf("Тренировки, на которые есть запись:")
	for _, t := range trainings {
		fmt.Printf("\nid: %d\nНазвание: %s\nДата и время: %s\n\n", t.ID, t.Name, t.DateTime)
	}

	return nil
}

func DeleteAssignment(a *registry.AppServiceFields, client *models.Client) {
	err := printassignment(a, client)
	if err != nil {
		return
	}

	var tID uint64
	fmt.Printf("Выберете тренировку и введите её id: ")
	_, err = fmt.Scanf("%d", &tID)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	err = a.ClientService.DeleteAssignment(ctx, client.ID, tID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Запись успешно удалена!")
}

const client_loop_string = `Меню клиента: 
	0 -- выйти
	1 -- посмотреть информацию о себе
	2 -- посмотреть расписание на неделю
	3 -- посмотреть расписание на выбраный промежуток времени 
	4 -- посмотреть тренеров по направлению
	5 -- записаться на тренировку
	6 -- отменить запись на тренировку
	7 -- посмотреть тренировки, на которые есть запись на неделю
Выберите действие: `

func clientMenu(a *registry.AppServiceFields, client *models.Client) {
	var num int = 1

	for num != 0 {
		fmt.Printf("\n\n%s", client_loop_string)

		_, err := fmt.Scanf("%d", &num)
		if err != nil {
			fmt.Printf("\nПункт меню введён некорректно!\n\n")
			continue
		}

		if num == 0 {
			return
		}

		switch num {
		case 0:
			return
		case 1:
			printInfoClient(client)
		case 2:
			printTrainingsOnWeek(a)
		case 3:
			printTrainings(a)
		case 4:
			printCoaches(a)
		case 5:
			CreateAssignment(a, client)
		case 6:
			DeleteAssignment(a, client)
		case 7:
			printassignment(a, client)
		default:
			fmt.Printf("\nНеверный пункт меню!\n\n")
		}
	}
}
