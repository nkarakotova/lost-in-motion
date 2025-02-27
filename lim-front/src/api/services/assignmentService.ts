import { AssignmentsApi, TrainingID } from '../client/api';
import { Configuration } from '../client/configuration';

class TrainingAssignmentService {
  private api: AssignmentsApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.api = new AssignmentsApi(configuration);
  }

  
  async createAssignment(trainingId: TrainingID, setError: (message: string) => void): Promise<void> {
    try {
      await this.api.createAssignment(trainingId);
      console.log('Запись на тренировку успешно создана!');
    } catch (error) {
      throw this.handleApiError(error, setError);
    }
  }

  async deleteAssignment(trainingId: TrainingID, setError: (message: string) => void): Promise<void> {
    try {
      await this.api.deleteAssignment(trainingId);
      console.log('Запись на тренировку успешно удалена!');
    } catch (error) {
      throw this.handleApiError(error, setError);
    }
  }

  private handleApiError(error: any, setError: (message: string) => void): string {
    if (error.response) {
      const status = error.response.status;
      let errorMessage = 'Неизвестная ошибка';

      switch (status) {
        case 401:
          errorMessage = 'Ошибка авторизации. Пожалуйста, войдите в систему.';
          break;
        case 403:
          errorMessage = 'У вас нет прав для выполнения этого действия.';
          break;
        case 404:
          errorMessage = 'Тренировка не найдена.';
          break;
        case 409:
          errorMessage = 'Вы уже записаны на тренировку в данное время.';
          break;
        case 422:
          errorMessage = 'Неверные параметры.';
          break;
        case 500:
          errorMessage = 'Внутренняя ошибка сервера.';
          break;
        default:
          errorMessage = `Неизвестная ошибка. Статус: ${status}`;
          break;
      }
      console.error(errorMessage);
      setError(errorMessage);  
      return errorMessage;
    } else if (error.request) {
      console.error('Ошибка при отправке запроса');
      setError('Ошибка при отправке запроса');
      return 'Ошибка при отправке запроса';
    } else {
      console.error('Неизвестная ошибка:', error.message);
      setError(`Неизвестная ошибка: ${error.message}`);
      return `Неизвестная ошибка: ${error.message}`;
    }
  }
}

export default TrainingAssignmentService;
