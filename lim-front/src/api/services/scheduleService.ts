import { ScheduleApi, Training } from '../client/api';
import { Configuration } from '../client/configuration';

export class ScheduleService {
  private scheduleApi: ScheduleApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.scheduleApi = new ScheduleApi(configuration);
  }

  public async getSchedule(): Promise<Training[]> {
    try {
      const response = await this.scheduleApi.getSchedule();

      const scheduleData: Training[] = response.data;
      console.log('Расписание успешно загружено');
      return scheduleData;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  private handleError(error: any): string {
    if (error.response) {
      const status = error.response.status;
      let errorMessage = 'Неизвестная ошибка';

      switch (status) {
        case 500:
          errorMessage = 'Ошибка сервера: невозможно получить расписание';
          break;
        default:
          errorMessage = `Неизвестная ошибка. Статус: ${status}`;
          break;
      }

      console.error(errorMessage);
      return errorMessage;
    } else if (error.request) {
      const requestErrorMessage = 'Ошибка при отправке запроса';
      console.error(requestErrorMessage);
      return requestErrorMessage;
    } else {
      const unknownErrorMessage = `Неизвестная ошибка: ${error.message}`;
      console.error(unknownErrorMessage);
      return unknownErrorMessage;
    }
  }
}
