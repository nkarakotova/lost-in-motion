import { HallsApi, HallInfo, Hall, GetHalls200ResponseInner } from '../client/api'; 
import { Configuration } from '../client/configuration'; 

export class HallService {
  private hallsApi: HallsApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }

    this.hallsApi = new HallsApi(configuration);
  }

  public async createHall(hallInfo: HallInfo, setError: (message: string) => void): Promise<void> {
    try {
      await this.hallsApi.createHall(hallInfo);
      console.log('Зал успешно создан');
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  public async getHalls(date: string, setError: (message: string) => void): Promise<{ time: string, halls: Hall[] }[]> {
    try {
      const response = await this.hallsApi.getHalls(date);

      const hallsByTime: { [key: string]: Hall[] } = {};

      response.data.forEach((item: GetHalls200ResponseInner) => {
        
        const time = item.timestamp ?? new Date().toISOString(); 
    
        
        const formattedTime = this.formatTime(time);

        
        if (!hallsByTime[formattedTime]) {
          hallsByTime[formattedTime] = [];
        }

        
        hallsByTime[formattedTime].push(...(item.hall || []));
      });

      return Object.entries(hallsByTime).map(([time, halls]) => ({
        time,
        halls,
      }));
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  private formatTime(timestamp: string): string {
    console.log('Полученный timestamp:', timestamp);
  
    const regex = /^time\.Date\((\d{4}),\s*time\.(\w+),\s*(\d{1,2}),\s*(\d{1,2}),\s*(\d{1,2}),\s*(\d{1,2}),\s*(\d{1,2}),\s*time\.(\w+)\)$/;
  
    const match = timestamp.match(regex);
    console.log('Результат матча:', match);
  
    if (match) {
      const year = match[1];
      const month = this.getMonthIndex(match[2]);
      const day = match[3];
      const hour = match[4];
      const minute = match[5];
      const second = match[6];
      const timezone = match[7];
  
      console.log('Извлеченные компоненты:', { year, month, day, hour, minute, second, timezone });
  
      const isoString = `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}T${hour}:${minute.padStart(2, '0')}:${second.padStart(2, '0')}Z`;
      console.log('ISO строка:', isoString);
  
      const date = new Date(isoString);
      if (isNaN(date.getTime())) {
        return 'Некорректное время';
      }
  
      const hours = date.getUTCHours().toString().padStart(2, '0');
      const minutes = date.getUTCMinutes().toString().padStart(2, '0');
  
      return `${hours}:${minutes}`;
    }
  
    const date = new Date(timestamp);
    if (isNaN(date.getTime())) {
      return 'Некорректное время';
    }
  
    const hours = date.getUTCHours().toString().padStart(2, '0');
    const minutes = date.getUTCMinutes().toString().padStart(2, '0');
  
    return `${hours}:${minutes}`;
  }

  private getMonthIndex(month: string): string {
    const months: { [key: string]: string } = {
      January: '01',
      February: '02',
      March: '03',
      April: '04',
      May: '05',
      June: '06',
      July: '07',
      August: '08',
      September: '09',
      October: '10',
      November: '11',
      December: '12',
    };
  
    if (months[month as keyof typeof months]) {
      return months[month as keyof typeof months];
    }
  
    return '01';
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
          errorMessage = 'Зал уже существует.';
          break;
        case 422:
          errorMessage = 'Неверные параметры запроса.';
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
