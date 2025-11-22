import React from 'react'
import { searchContents, getMetadata, SearchResponse, ContentTypeMetadata, SortOptionMetadata } from '../api/search'
import { SearchForm } from '../components/SearchForm'
import { ContentTable } from '../components/ContentTable'
import { Pagination } from '../components/Pagination'

export function SearchPage() {
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [searchResult, setSearchResult] = React.useState<SearchResponse | null>(null)
  const [contentTypes, setContentTypes] = React.useState<ContentTypeMetadata[]>([])
  const [sortOptions, setSortOptions] = React.useState<SortOptionMetadata[]>([])
  const [searchParams, setSearchParams] = React.useState({
    query: '',
    type: 'all',
    sort: '',
    page: 1,
    page_size: 10,
  })

  // Fetch metadata once on mount
  React.useEffect(() => {
    const fetchMetadata = async () => {
      try {
        const metadata = await getMetadata()
        setContentTypes(metadata.content_types)
        setSortOptions(metadata.sort_options)
        
        if (metadata.sort_options.length > 0) {
          setSearchParams((prev) => ({ 
            ...prev, 
            sort: metadata.sort_options[0].id 
          }))
        }
        
        if (metadata.pagination) {
          setSearchParams((prev) => ({ 
            ...prev, 
            page_size: metadata.pagination.default_page_size 
          }))
        }
      } catch (err) {
        console.error('Failed to fetch metadata:', err)
        setError('Failed to load configuration')
      }
    }
    
    fetchMetadata()
  }, []) // Empty dependency - only run once

  // Perform search when params change
  React.useEffect(() => {
    const performSearch = async () => {
      if (!searchParams.sort) return // Wait for sort to be set

      setLoading(true)
      setError(null)

      try {
        const result = await searchContents({
          ...searchParams,
          type: searchParams.type === 'all' ? undefined : searchParams.type,
        })
        setSearchResult(result)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred')
        setSearchResult(null)
      } finally {
        setLoading(false)
      }
    }

    performSearch()
  }, [searchParams]) // Only depend on searchParams

  const handleSearch = React.useCallback((params: { query: string; type: string; sort: string }) => {
    setSearchParams(prev => ({ ...prev, ...params, page: 1 }))
  }, [])

  const handlePageChange = React.useCallback((page: number) => {
    setSearchParams(prev => ({ ...prev, page }))
  }, [])

  const totalPages = searchResult
    ? Math.ceil(searchResult.total / searchParams.page_size)
    : 0

  return (
    <div className="app-shell">
      <div className="frame">
        <header className="hero">
          <h1 className="hero-title">Content Search</h1>
          <p className="hero-subtitle">Search and discover content from multiple sources</p>
        </header>

        <SearchForm
          contentTypes={contentTypes}
          sortOptions={sortOptions}
          onSearch={handleSearch}
        />

        {error && <div className="alert">{error}</div>}

        <section className="panel">
          <ContentTable
            items={searchResult?.items || []}
            loading={loading}
            total={searchResult?.total || 0}
          />

          {searchResult && searchResult.total > 0 && (
            <Pagination
              currentPage={searchParams.page}
              totalPages={totalPages}
              onPageChange={handlePageChange}
            />
          )}
        </section>
      </div>
    </div>
  )
}
