import { AuthApi } from '../client/api';
import { ClientInfo, ClientLoginInfo, Token } from '../client/api';
import { Configuration } from '../client/configuration';

export class AuthService {
  private authApi: AuthApi;

  constructor(configuration?: Configuration) {
    if (configuration) {
      configuration.addAuthorizationHeader();
    }
    this.authApi = new AuthApi(configuration);
  }

  public async signup(clientInfo: ClientInfo, setError: (message: string) => void): Promise<void> {
    try {
      await this.authApi.signup(clientInfo);
      console.log('Клиент успешно зарегистрирован');
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  public async login(clientLoginInfo: ClientLoginInfo, setError: (message: string) => void): Promise<Token> {
    try {
      const response = await this.authApi.login(clientLoginInfo);
      const tokenData: Token = response.data;

      if (tokenData) {
        localStorage.setItem('jwt_token', JSON.stringify(tokenData).replace(/^"|"$/g, ''));
        console.log('Токен успешно сохранен:', tokenData);
      } else {
        console.error('Токен отсутствует в ответе');
        setError('Не удалось получить токен');
        throw new Error('Не удалось получить токен');
      }

      console.log('Клиент успешно вошел');
      return tokenData;
    } catch (error) {
      throw this.handleError(error, setError);
    }
  }

  public getToken(): string | null {
    return localStorage.getItem('jwt_token');
  }

  public logout(): void {
    localStorage.removeItem('jwt_token');
    console.log('Клиент успешно вышел из системы');
  }

  private handleError(error: any, setError: (message: string) => void): string {
    if (error.response) {
      const status = error.response.status;
      let errorMessage = 'Неизвестная ошибка';

      switch (status) {
        case 401:
          errorMessage = 'Неверный пароль.';
          break;
        case 404:
          errorMessage = 'Клиент не найден.';
          break;
        case 409:
          errorMessage = 'Клиент уже существует.';
          break;
        case 422:
          errorMessage = 'Неверные параметры.';
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
      const requestErrorMessage = 'Ошибка при отправке запроса';
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
