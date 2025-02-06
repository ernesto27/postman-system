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
  host: string[] | string;
  port: string;
  path: string[];
  query?: QueryParam[];
  variable?: UrlVariable[];
}

export interface UrlVariable {
  key: string;
  value: string;
  description?: string;
}

export interface QueryParam {
  key: string;
  value: string;
  disabled?: boolean;
}

export interface PostmanAuth {
  type: 'bearer' | 'apikey' | 'oauth2';
  bearer?: Array<{
    key: string;
    value: string;
    type: string;
  }>;
  apikey?: Array<{
    key: string;
    value: string | boolean;
    type: string;
  }>;
}

export interface PostmanRequest {
  auth?: PostmanAuth;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  header: HeaderItem[];
  body?: PostmanRequestBody;
  url: PostmanUrl | string;
  description?: string;
}

export interface HeaderItem {
  key: string;
  value: string;
  disabled?: boolean;
}

export interface PostmanItem {
  id?: string;
  name: string;
  request?: PostmanRequest;
  response?: any[];
  item?: PostmanItem[];
  description?: string;
}

export interface PostmanCollection {
  id: string;
  name: string;
  info: {
    _postman_id: string;
    name: string;
    description?: string;
    schema: string;
    _exporter_id?: string;
    _collection_link?: string;
  };
  item: PostmanItem[];
  variable?: Variable[];
  createdAt?: string;
  updatedAt?: string;
}

export interface Variable {
  key: string;
  value: string;
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

export interface PostmanResponse {
  name: string;
  originalRequest: PostmanRequest;
  status?: string;
  code: number;
  _postman_previewlanguage?: string;
  header: HeaderItem[];
  cookie: any[];
  body: any;
}
