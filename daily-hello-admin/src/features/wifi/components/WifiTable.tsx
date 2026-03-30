import Table from '@/shared/components/Table'
import type { Wifi } from '../types'

type Props = {
  data: Wifi[]
  loading: boolean
  onDelete: (id: number) => void
}

export default function WifiTable({ data, loading, onDelete }: Props) {
  const columns = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: 'Name' },
    { key: 'code', title: 'Code' },
    { key: 'ssid', title: 'SSID' },
    { key: 'bssid', title: 'BSSID' },
    {
      key: 'branch_name',
      title: 'Branch Name',
      render: (item: Wifi) => item.branch_name || item.branch?.name || '',
    },
    {
      key: 'actions',
      title: '',
      render: (item: Wifi) => (
        <button onClick={() => onDelete(item.id)}>Xoa</button>
      ),
    },
  ]

  return <Table columns={columns} data={data} loading={loading} />
}
