import { fetchJSON } from './client'

// Match protobuf ContentItem message
export interface ContentItem {
  id: number
  title: string
  content_type: string
  score: number
  published_at: string
  provider_name: string
}

export interface SearchResponse {
  items: ContentItem[]
  page: number
  page_size: number
  total: number
}

export interface SearchParams {
  query?: string
  type?: string
  sort?: string
  page?: number
  page_size?: number
}

export interface ContentTypeMetadata {
  id: string
  display_name: string
}

export interface SortOptionMetadata {
  id: string
  display_name: string
}

export interface PaginationMetadata {
  default_page_size: number
  max_page_size: number
}

export interface MetadataResponse {
  content_types: ContentTypeMetadata[]
  sort_options: SortOptionMetadata[]
  pagination: PaginationMetadata
}

export async function getMetadata(): Promise<MetadataResponse> {
  return fetchJSON<MetadataResponse>('/api/v1/metadata')
}

export async function searchContents(params: SearchParams): Promise<SearchResponse> {
  const searchParams = new URLSearchParams()
  
  if (params.query) searchParams.set('query', params.query)
  if (params.type) searchParams.set('type', params.type)
  if (params.sort) searchParams.set('sort', params.sort)
  if (params.page) searchParams.set('page', params.page.toString())
  if (params.page_size) searchParams.set('page_size', params.page_size.toString())
  
  const url = `/api/v1/search?${searchParams.toString()}`
  return fetchJSON<SearchResponse>(url)
}

export async function getContent(id: number): Promise<{ content: ContentItem }> {
  return fetchJSON<{ content: ContentItem }>(`/api/v1/contents/${id}`)
}
