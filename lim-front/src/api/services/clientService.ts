import { ClientsApi, ClientChangePassword, Training } from '../client/api';
import { Configuration } from '../client/configuration';

export class ClientService {
  private clientsApi: ClientsApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.clientsApi = new ClientsApi(configuration);
  }

  public async changePassword(newPassword: string, setError: (message: string) => void): Promise<void> {
    const clientPasswordData: ClientChangePassword = { password: newPassword };

    try {
      await this.clientsApi.changePassword(clientPasswordData);
      console.log('Пароль успешно изменен');
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  public async getTrainingsByClient(setError: (message: string) => void): Promise<Training[]> {
    try {
      const response = await this.clientsApi.getTrainingsByClient();
      const trainingsData: Training[] = response.data;

      console.log('Тренировки клиента успешно загружены');
      return trainingsData;
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
        case 404:
          errorMessage = 'Клиент или данные не найдены.';
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
