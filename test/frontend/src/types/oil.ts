export interface OilPrice {
  id: number;
  date: string;
  oil_type: string;
  price: number;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message: string;
}
