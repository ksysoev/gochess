interface APISettings {
    token: string;
    headers: Headers;
    baseURL: string;
}

const APIConfig: APISettings = {
  token: '',
  headers: new Headers({
    Accept: 'application/json',
  }),
  baseURL: 'http://localhost:8081',
};

export { APISettings, APIConfig };
