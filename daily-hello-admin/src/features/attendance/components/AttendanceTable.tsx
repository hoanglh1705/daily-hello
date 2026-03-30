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
    {
      key: 'user',
      title: 'Nhan vien',
      render: (item: Attendance) => item.user.name,
    },
    {
      key: 'branch',
      title: 'Chi nhanh',
      render: (item: Attendance) => item.branch.name,
    },
    {
      key: 'check_in_time',
      title: 'Check in',
      render: (item: Attendance) => formatDateTime(item.check_in_time),
    },
    {
      key: 'check_out_time',
      title: 'Check out',
      render: (item: Attendance) =>
        item.check_out_time ? formatDateTime(item.check_out_time) : '-',
    },
    { key: 'status', title: 'Trang thai' },
  ]

  return <Table columns={columns} data={data} loading={loading} />
}
