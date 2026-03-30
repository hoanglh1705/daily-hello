import type { User } from '../types'
import Table from '@/shared/components/Table'

type Props = {
  data: User[]
  loading: boolean
  onEdit?: (user: User) => void
}

function getInitials(name: string) {
  return name
    .split(' ')
    .map((w) => w[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

function StatusBadge({ status }: { status: string }) {
  const s = status.toLowerCase()
  const classMap: Record<string, string> = {
    active: 'active',
    inactive: 'inactive',
    remote: 'remote',
    on_leave: 'on-leave',
    pending: 'pending',
  }
  return (
    <span className={`status-badge ${classMap[s] || 'inactive'}`}>
      {status.replace('_', ' ')}
    </span>
  )
}

export default function UserTable({ data, loading, onEdit }: Props) {
  const columns = [
    {
      key: 'name',
      title: 'Tên',
      render: (u: User) => (
        <div className="cell-user">
          <div className="cell-avatar">{getInitials(u.name)}</div>
          <div className="cell-user-info">
            <span className="cell-user-name">{u.name}</span>
            <span className="cell-user-sub">{u.email}</span>
          </div>
        </div>
      ),
    },
    { key: 'code', title: 'Mã NV' },
    { key: 'phone', title: 'SĐT' },
    {
      key: 'branch',
      title: 'Chi nhánh',
      render: (u: User) => u.branch?.name || '—',
    },
    { key: 'role', title: 'Vai trò' },
    {
      key: 'status',
      title: 'Trạng thái',
      render: (u: User) => <StatusBadge status={u.status} />,
    },
    {
      key: 'action',
      title: 'Thao tác',
      width: '100px',
      render: (u: User) => (
        <div className="cell-actions">
          <button
            className="action-btn action-btn-primary"
            title="Chỉnh sửa"
            onClick={() => onEdit?.(u)}
            disabled={!onEdit}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
              <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
            </svg>
          </button>
        </div>
      ),
    },
  ]

  return <Table data={data} columns={columns} loading={loading} />
}
