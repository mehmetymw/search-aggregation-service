import React from 'react'

interface ContentTypeMetadata {
  id: string
  display_name: string
}

interface SortOptionMetadata {
  id: string
  display_name: string
}

interface Props {
  contentTypes: ContentTypeMetadata[]
  sortOptions: SortOptionMetadata[]
  onSearch: (params: { query: string; type: string; sort: string }) => void
}

export function SearchForm({ contentTypes, sortOptions, onSearch }: Props) {
  const [query, setQuery] = React.useState('')
  const [type, setType] = React.useState('all')
  const [sort, setSort] = React.useState('')

  // Set default sort only once when options load
  React.useEffect(() => {
    if (sortOptions.length > 0 && !sort) {
      setSort(sortOptions[0].id)
    }
  }, [sortOptions, sort])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSearch({ query, type, sort })
  }

  // Debounced search - only trigger on user input
  React.useEffect(() => {
    if (!sort) return // Don't search until sort is set
    
    const timeoutId = setTimeout(() => {
      onSearch({ query, type, sort })
    }, 300)

    return () => clearTimeout(timeoutId)
  }, [query, type, sort]) // Removed onSearch from dependencies to prevent loop

  return (
    <form onSubmit={handleSubmit} className="search-form">
      <div className="form-group" style={{ flex: 2 }}>
        <label className="label" htmlFor="query">
          Search
        </label>
        <input
          id="query"
          type="text"
          className="input"
          placeholder="Search for content..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      <div className="form-group">
        <label className="label" htmlFor="type">
          Type
        </label>
        <select
          id="type"
          className="select"
          value={type}
          onChange={(e) => setType(e.target.value)}
        >
          {contentTypes.map((ct) => (
            <option key={ct.id} value={ct.id}>
              {ct.display_name}
            </option>
          ))}
        </select>
      </div>

      <div className="form-group">
        <label className="label" htmlFor="sort">
          Sort
        </label>
        <select
          id="sort"
          className="select"
          value={sort}
          onChange={(e) => setSort(e.target.value)}
        >
          {sortOptions.map((so) => (
            <option key={so.id} value={so.id}>
              {so.display_name}
            </option>
          ))}
        </select>
      </div>

      <div className="form-actions">
        <button type="submit" className="button">
          Search
        </button>
      </div>
    </form>
  )
}
