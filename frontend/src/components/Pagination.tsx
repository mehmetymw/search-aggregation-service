interface Props {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
}

export function Pagination({ currentPage, totalPages, onPageChange }: Props) {
  const canGoPrev = currentPage > 1
  const canGoNext = currentPage < totalPages

  return (
    <div className="pagination">
      <button
        className="page-btn"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={!canGoPrev}
        aria-label="Previous page"
      >
        ← Previous
      </button>

      <div className="page-info">
        Page <strong>{currentPage}</strong> of <strong>{totalPages}</strong>
      </div>

      <button
        className="page-btn"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={!canGoNext}
        aria-label="Next page"
      >
        Next →
      </button>
    </div>
  )
}
