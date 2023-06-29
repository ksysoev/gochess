import { APIConfig } from '@/api/config';

class APIClient {
    private baseURL: string;

    private headers: Headers;

    private token: string;

    private static instance: APIClient;

    private eventSource: EventSource|null = null;

    constructor() {
      this.baseURL = APIConfig.baseURL;
      this.headers = APIConfig.headers;
      this.token = APIConfig.token;
    }

    public static getInstance(): APIClient {
      if (!APIClient.instance) {
        APIClient.instance = new APIClient();
      }

      return APIClient.instance;
    }

    private async get(path: string) : Promise<string> {
      const response = await fetch(this.baseURL + path, {
        method: 'GET',
        headers: this.headers,
      });
      return response.text();
    }

    private async post(path: string, body: string) : Promise<string> {
      const response = await fetch(this.baseURL + path, {
        method: 'POST',
        headers: this.headers,
        body,
      });
      return response.text();
    }

    public async getGame(id: string): Promise<any> {
      const response = await this.get(`/game/${id}`);

      return JSON.parse(response);
    }

    public async findMatch(playerName: string) {
      return this.post('/match', JSON.stringify({
        name: playerName,
      }));
    }

    public async listen(eventName: string, callback: (event: Event) => void) {
      if (this.eventSource === null) {
        this.eventSource = new EventSource(`${this.baseURL}/notifier`);
      }

      this.eventSource.addEventListener(eventName, callback);

      await new Promise<void>((resolve) => {
        if (this.eventSource === null) {
          resolve();
          return;
        }

        if (this.eventSource.readyState === EventSource.OPEN) {
          resolve();
          return;
        }
        this.eventSource.addEventListener('open', () => {
          resolve();
        });
      });
    }

    public async forget(eventName: string, callback: (event: Event) => void) {
      if (this.eventSource === null) {
        return;
      }

      this.eventSource.removeEventListener(eventName, callback);
    }
}

interface EventGameStarted {
    GameID: string;
    PlayerBlack: string;
    PlayerWhite: string;
    Position: string;
  }

export { APIClient, EventGameStarted };
