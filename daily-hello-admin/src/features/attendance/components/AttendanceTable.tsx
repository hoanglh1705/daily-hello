import Table from '@/shared/components/Table'
import { formatDateTime } from '@/shared/utils/formatDate'
import type { Attendance } from '../types'

type Props = {
  data: Attendance[]
  loading: boolean
}

export default function AttendanceTable({ data, loading }: Props) {
  const columns = [
    { key: 'id', title: 'ID' },
    { key: 'user_name', title: 'Nhan vien' },
    { key: 'branch_name', title: 'Chi nhanh' },
    {
      key: 'check_in_at',
      title: 'Check in',
      render: (item: Attendance) => formatDateTime(item.check_in_at),
    },
    {
      key: 'check_out_at',
      title: 'Check out',
      render: (item: Attendance) =>
        item.check_out_at ? formatDateTime(item.check_out_at) : '-',
    },
    { key: 'status', title: 'Trang thai' },
  ]

  return <Table columns={columns} data={data} loading={loading} />
}
