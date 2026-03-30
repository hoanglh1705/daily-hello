import Table from '@/shared/components/Table'
import { formatDate } from '@/shared/utils/formatDate'
import type { Branch } from '../types'

type Props = {
  data: Branch[]
  loading: boolean
  onEdit: (branch: Branch) => void
  onDelete: (id: number) => void
}

function StatusBadge({ status }: { status: string }) {
  const s = status.toLowerCase()
  const classMap: Record<string, string> = {
    active: 'active',
    inactive: 'inactive',
  }
  return (
    <span className={`status-badge ${classMap[s] || 'inactive'}`}>
      {status}
    </span>
  )
}

export default function BranchTable({ data, loading, onEdit, onDelete }: Props) {
  const columns = [
    {
      key: 'name',
      title: 'Tên chi nhánh',
      render: (item: Branch) => (
        <div className="cell-user">
          <div className="cell-avatar" style={{ background: 'linear-gradient(135deg, #dbeafe, #bfdbfe)', color: '#1e40af' }}>
            {item.name.charAt(0).toUpperCase()}
          </div>
          <div className="cell-user-info">
            <span className="cell-user-name">{item.name}</span>
            <span className="cell-user-sub">{item.branch_code}</span>
          </div>
        </div>
      ),
    },
    { key: 'address', title: 'Địa chỉ' },
    { key: 'parent_branch_code', title: 'Chi nhánh cha' },
    {
      key: 'status',
      title: 'Trạng thái',
      render: (item: Branch) => <StatusBadge status={item.status} />,
    },
    {
      key: 'created_at',
      title: 'Ngày tạo',
      render: (item: Branch) => (
        <span style={{ color: '#6b7280', fontSize: '13px' }}>{formatDate(item.created_at)}</span>
      ),
    },
    {
      key: 'actions',
      title: 'Thao tác',
      width: '120px',
      render: (item: Branch) => (
        <div className="cell-actions">
          <button className="action-btn action-btn-primary" title="Sửa" onClick={() => onEdit(item)}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
              <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
            </svg>
          </button>
          <button className="action-btn action-btn-danger" title="Xóa" onClick={() => onDelete(item.id)}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="3 6 5 6 21 6" />
              <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
            </svg>
          </button>
        </div>
      ),
    },
  ]

  return <Table columns={columns} data={data} loading={loading} />
}
