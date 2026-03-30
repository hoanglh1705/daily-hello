type Props = {
  page: number
  limit: number
  total: number
  onPageChange: (page: number) => void
}

export default function Pagination({ page, limit, total, onPageChange }: Props) {
  const totalPages = Math.ceil(total / limit)

  if (totalPages <= 1) return null

  return (
    <div className="pagination">
      <button disabled={page <= 1} onClick={() => onPageChange(page - 1)}>
        Prev
      </button>
      <span>
        {page} / {totalPages}
      </span>
      <button disabled={page >= totalPages} onClick={() => onPageChange(page + 1)}>
        Next
      </button>
    </div>
  )
}
