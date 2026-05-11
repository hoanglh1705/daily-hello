import type { Holiday } from '../types'

type Props = {
  data: Holiday[]
  loading: boolean
  onEdit: (holiday: Holiday) => void
  onDelete: (id: number) => void
}

export default function HolidayTable({ data, loading, onEdit, onDelete }: Props) {
  if (loading) {
    return <div className="table-loading">Đang tải dữ liệu...</div>
  }

  return (
    <div className="table-container">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Tên ngày nghỉ</th>
            <th>Ngày</th>
            <th>Mô tả</th>
            <th>Thao tác</th>
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={5} style={{ textAlign: 'center', padding: '2rem' }}>
                Chưa có ngày nghỉ nào
              </td>
            </tr>
          ) : (
            data.map((item) => (
              <tr key={item.id}>
                <td>{item.id}</td>
                <td>
                  <strong>{item.name}</strong>
                </td>
                <td>
                  <span className="badge badge-info">
                    {new Date(item.date).toLocaleDateString('vi-VN')}
                  </span>
                </td>
                <td>{item.description || '—'}</td>
                <td>
                  <div className="action-btns">
                    <button className="btn-icon btn-edit" title="Sửa" onClick={() => onEdit(item)}>
                      ✏️
                    </button>
                    <button className="btn-icon btn-delete" title="Xóa" onClick={() => onDelete(item.id)}>
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}
