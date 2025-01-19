export interface PostmanRequestBody {
  mode: 'raw' | 'formdata';
  raw?: string;
  formdata?: FormDataItem[];
  options?: {
    raw: {
      language: string;
    };
  };
}

export interface FormDataItem {
  key: string;
  value?: string;
  type: 'text' | 'file';
  src?: string;
  disabled?: boolean;
}

export interface PostmanUrl {
  raw: string;
  protocol: string;
  host: string[];
  port: string;
  path: string[];
  query?: QueryParam[];
}

export interface QueryParam {
  key: string;
  value: string;
  disabled?: boolean;
}

export interface PostmanAuth {
  type: 'bearer';
  bearer: Array<{
    key: string;
    value: string;
    type: string;
  }>;
}

export interface PostmanRequest {
  auth?: PostmanAuth;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  header: HeaderItem[];
  body?: PostmanRequestBody;
  url: PostmanUrl;
}

export interface HeaderItem {
  key: string;
  value: string;
  disabled?: boolean;
}

export interface PostmanItem {
  id: string;
  name: string;
  request: PostmanRequest;
  response: any[];
}

export interface PostmanCollection {
  id: string;
  name: string;
  items: PostmanItem[];
  createdAt?: string;
  updatedAt?: string;
}

export interface TabState {
  id: number;
  title: string;
  request: PostmanRequest;
  response?: {
    data: string;
    status: string;
    headers: Record<string, string>;
    time: number;
    size: number;
  };
}

export interface LocalStorage {
  collections: PostmanCollection[];
  tabs: TabState[];
  activeTabId: number;
  history: PostmanHistoryItem[];
}

export interface PostmanHistoryItem {
  id: string;
  method: string;
  url: string;
  timestamp: string;
  request: PostmanRequest;
  response?: {
    data: string;
    status: string;
    headers: Record<string, string>;
    time: number;
    size: number;
  };
}
