import Table from '@/shared/components/Table'
import { formatDate } from '@/shared/utils/formatDate'
import type { Branch } from '../types'

type Props = {
  data: Branch[]
  loading: boolean
  onEdit: (branch: Branch) => void
  onDelete: (id: number) => void
}

export default function BranchTable({ data, loading, onEdit, onDelete }: Props) {
  const columns = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: 'Ten chi nhanh' },
    { key: 'address', title: 'Dia chi' },
    {
      key: 'created_at',
      title: 'Ngay tao',
      render: (item: Branch) => formatDate(item.created_at),
    },
    {
      key: 'actions',
      title: '',
      render: (item: Branch) => (
        <div>
          <button onClick={() => onEdit(item)}>Sua</button>
          <button onClick={() => onDelete(item.id)}>Xoa</button>
        </div>
      ),
    },
  ]

  return <Table columns={columns} data={data} loading={loading} />
}
