type Column<T> = {
  key: string
  title: string
  width?: string
  render?: (item: T) => React.ReactNode
}

type Props<T> = {
  columns: Column<T>[]
  data: T[]
  loading?: boolean
  getRowClassName?: (item: T, index: number) => string | undefined
}

export type { Column }

export default function Table<T extends { id: string | number }>({
  columns,
  data,
  loading,
  getRowClassName,
}: Props<T>) {
  if (loading) {
    return (
      <div className="table-loading">
        <div className="table-loading-spinner" />
        <span>Đang tải dữ liệu...</span>
      </div>
    )
  }

  return (
    <div className="table-wrapper">
      <table className="shared-table">
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={col.key} style={col.width ? { width: col.width } : undefined}>
                {col.title}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} className="table-empty">
                Không có dữ liệu
              </td>
            </tr>
          ) : (
            data.map((item, index) => (
              <tr key={item.id} className={getRowClassName?.(item, index)}>
                {columns.map((col) => (
                  <td key={col.key}>
                    {col.render
                      ? col.render(item)
                      : (item as Record<string, unknown>)[col.key] as React.ReactNode}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}
