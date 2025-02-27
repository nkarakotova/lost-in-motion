import { TrainingsApi, TrainingInfo, TrainingID } from '../client/api';
import { Configuration } from '../client/configuration';

export class TrainingService {
  private trainingsApi: TrainingsApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.trainingsApi = new TrainingsApi(configuration);
  }

  public async createTraining(trainingInfo: TrainingInfo, setError: (message: string) => void): Promise<void> {
    try {
      console.log(trainingInfo);
      await this.trainingsApi.createTraining(trainingInfo);
      console.log('Тренировка успешно создана');
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  public async deleteTraining(trainingId: TrainingID, setError: (message: string) => void): Promise<void> {
    try {
      await this.trainingsApi.deleteTraining(trainingId);
      console.log('Тренировка успешно удалена');
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
          errorMessage = 'Ошибка авторизации: необходимо войти в систему';
          break;
        case 403:
          errorMessage = 'Нет прав для выполнения операции';
          break;
        case 404:
          errorMessage = 'Тренировка не найдена';
          break;
        case 409:
          errorMessage = 'У выбранного тренера уже есть тренировка в данное время';
          break;
        case 422:
          errorMessage = 'Некорректные параметры';
          break;
        case 500:
          errorMessage = 'Ошибка сервера';
          break;
        default:
          errorMessage = 'Неизвестная ошибка';
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
