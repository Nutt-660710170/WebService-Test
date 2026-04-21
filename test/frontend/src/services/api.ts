import type { ApiResponse, OilPrice } from '../types/oil';

export const fetchPrices = async (oilType?: string): Promise<OilPrice[]> => {
  const url = new URL('/api/v1/list', window.location.origin);
  if (oilType) {
    url.searchParams.append('oil_type', oilType);
  }
  
  const response = await fetch(url.toString());
  if (!response.ok) {
    throw new Error(`Error fetching prices: ${response.statusText}`);
  }
  const json: ApiResponse<OilPrice[]> = await response.json();
  return json.data;
};

export const pullPrices = async (): Promise<OilPrice[]> => {
  const response = await fetch('/api/v1/pull', { method: 'POST' });
  if (!response.ok) {
    throw new Error(`Error pulling prices: ${response.statusText}`);
  }
  const json: ApiResponse<OilPrice[]> = await response.json();
  return json.data;
};
