export interface ShortenResponse {
  short_url: string;
  code: string;
}

export interface UrlStats {
  short_code: string;
  original: string;
  clicks: number;
  created_at: string;
}