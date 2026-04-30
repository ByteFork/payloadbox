export interface Request {
  method: string;
  path: string;
  query?: string;
  headers: Record<string, string[]>;
  body: string;
  remote_addr: string;
  host: string;
  content_length: number;
}

export interface Response {
  status_code: number;
  status_text?: string;
  headers: Record<string, string[]>;
  body: string;
  size_in_bytes: number;
}

export interface RequestRecord {
  id: string;
  created_at: string;
  request: Request;
  response: Response;
  duration_ns: number;
}

export type PageType = "requests" | "docs" | "settings";

export type MethodFilter = "ALL" | "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
