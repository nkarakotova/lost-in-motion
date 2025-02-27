import { CoachesApi, Coach, CoachInfo } from '../client/api';
import { Configuration } from '../client/configuration';

export class CoachService {
  private coachesApi: CoachesApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.coachesApi = new CoachesApi(configuration);
  }

  public async createCoach(coachInfo: CoachInfo, setError: (message: string) => void): Promise<void> {
    try {
      await this.coachesApi.createCoach(coachInfo);
      console.log('Тренер успешно создан');
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  
  public async getCoaches(setError: (message: string) => void): Promise<Coach[]> {
    try {
      const response = await this.coachesApi.getCoaches();
      const coaches = response.data;
      console.log('Список тренеров:', coaches);
      return coaches;
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  private handleError(error: any, setError: (message: string) => void): string {
    if (error.response) {
      const status = error.response.status;
      let errorMessage = 'Неизвестная ошибка';

      switch (status) {
        case 401:
          errorMessage = 'Ошибка авторизации: необходимо войти в систему.';
          break;
        case 403:
          errorMessage = 'Нет прав для выполнения операции.';
          break;
        case 409:
          errorMessage = 'Тренер уже существует.';
          break;
        case 500:
          errorMessage = 'Ошибка сервера.';
          break;
        default:
          errorMessage = `Неизвестная ошибка. Статус: ${status}`;
          break;
      }

      console.error(errorMessage);
      setError(errorMessage); 
      return errorMessage;
    } else if (error.request) {
      const requestErrorMessage = 'Ошибка при отправке запроса.';
      console.error(requestErrorMessage);
      setError(requestErrorMessage);
      return requestErrorMessage;
    } else {
      const unknownErrorMessage = `Неизвестная ошибка: ${error.message}`;
      console.error(unknownErrorMessage);
      setError(unknownErrorMessage);
      return unknownErrorMessage;
    }
  }
}
