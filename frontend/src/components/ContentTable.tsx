interface ContentItem {
  id: number
  title: string
  content_type: string
  score: number
  published_at: string
  provider_name: string
}

interface Props {
  items: ContentItem[]
  loading: boolean
  total: string | number
}

const SearchIconLarge = () => (
  <svg 
    width="64" 
    height="64" 
    viewBox="0 0 24 24" 
    fill="none" 
    stroke="currentColor" 
    strokeWidth="2" 
    strokeLinecap="round" 
    strokeLinejoin="round"
    style={{ opacity: 0.2, marginBottom: '1rem' }}
  >
    <circle cx="11" cy="11" r="8"></circle>
    <path d="m21 21-4.35-4.35"></path>
  </svg>
)

export function ContentTable({ items, loading, total }: Props) {
  if (loading) {
    return <div className="loading"></div>
  }

  if (items.length === 0) {
    return (
      <div className="empty-state">
        <SearchIconLarge />
        <div>No content found. Try adjusting your search criteria.</div>
      </div>
    )
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  return (
    <div className="table-wrapper">
      <table className="table">
        <thead>
          <tr className="tr">
            <th className="th">Title</th>
            <th className="th">Type</th>
            <th className="th">Published</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.id} className="tr">
              <td className="td">
                <div style={{ fontWeight: 500 }}>
                  {item.title}
                </div>
              </td>
               <td className="td">
                <div style={{ fontWeight: 500 }}>
                  {item.content_type}
                </div>
              </td>
              <td className="td">
                <div style={{ color: 'var(--text-secondary)' }}>
                  {formatDate(item.published_at)}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      
      <div style={{ 
        padding: '1rem 1.5rem', 
        borderTop: '1px solid var(--border)',
        fontSize: '0.875rem',
        color: 'var(--text-secondary)',
        textAlign: 'center'
      }}>
        Showing {items.length} of {total} results
      </div>
    </div>
  )
}
