type Props = {
  page: number
  limit: number
  total: number
  onPageChange: (page: number) => void
}

export default function Pagination({ page, limit, total, onPageChange }: Props) {
  const totalPages = Math.ceil(total / limit)

  if (totalPages <= 1 && total === 0) return null

  const from = (page - 1) * limit + 1
  const to = Math.min(page * limit, total)

  const getPageNumbers = () => {
    const pages: (number | '...')[] = []
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) pages.push(i)
    } else {
      pages.push(1)
      if (page > 3) pages.push('...')
      for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) {
        pages.push(i)
      }
      if (page < totalPages - 2) pages.push('...')
      pages.push(totalPages)
    }
    return pages
  }

  return (
    <div className="pagination">
      <span className="pagination-info">
        Hiển thị {from} đến {to} của {total} bản ghi
      </span>
      <div className="pagination-controls">
        <button
          className="pagination-btn"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          Previous
        </button>
        {getPageNumbers().map((p, i) =>
          p === '...' ? (
            <span key={`dots-${i}`} className="pagination-dots">...</span>
          ) : (
            <button
              key={p}
              className={`pagination-page ${page === p ? 'active' : ''}`}
              onClick={() => onPageChange(p)}
            >
              {p}
            </button>
          )
        )}
        <button
          className="pagination-btn"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  )
}
