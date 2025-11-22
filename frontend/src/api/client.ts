const BASE_URL = import.meta?.env?.VITE_API_URL || 'http://localhost:8081'

export async function fetchJSON<T>(url: string): Promise<T> {
  const response = await fetch(`${BASE_URL}${url}`)
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  return response.json()
}
