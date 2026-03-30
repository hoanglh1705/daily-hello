import type { User } from '../types'
import Table from '@/shared/components/Table'

type Props = {
  data: User[]
  loading: boolean
  onEdit?: (user: User) => void
}

export default function UserTable({ data, loading, onEdit }: Props) {
  const columns = [
    { key: 'id', title: 'ID' },
    { key: 'code', title: 'Code' },
    { key: 'name', title: 'Ten' },
    { key: 'email', title: 'Email' },
    { key: 'phone', title: 'SDT' },
    { key: 'role', title: 'Role' },
    { key: 'status', title: 'Status' },
    {
      key: 'branch',
      title: 'Chi nhanh',
      render: (u: User) => u.branch?.name || '',
    },
    {
      key: 'action',
      title: 'Thao tác',
      render: (u: User) => (
        <button onClick={() => onEdit?.(u)} disabled={!onEdit}>
          Sửa
        </button>
      ),
    },
  ]

  return <Table data={data} columns={columns} loading={loading} />
}
