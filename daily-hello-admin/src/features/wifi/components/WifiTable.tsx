import Table from '@/shared/components/Table'
import type { Wifi } from '../types'

type Props = {
  data: Wifi[]
  loading: boolean
  onDelete: (id: number) => void
}

export default function WifiTable({ data, loading, onDelete }: Props) {
  const columns = [
    {
      key: 'name',
      title: 'Tên WiFi',
      render: (item: Wifi) => (
        <div className="cell-user">
          <div className="cell-avatar" style={{ background: 'linear-gradient(135deg, #fef3c7, #fde68a)', color: '#92400e' }}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" width="18" height="18">
              <path d="M5 12.55a11 11 0 0114.08 0" />
              <path d="M1.42 9a16 16 0 0121.16 0" />
              <path d="M8.53 16.11a6 6 0 016.95 0" />
              <circle cx="12" cy="20" r="1" fill="currentColor" />
            </svg>
          </div>
          <div className="cell-user-info">
            <span className="cell-user-name">{item.name}</span>
            <span className="cell-user-sub">{item.code}</span>
          </div>
        </div>
      ),
    },
    { key: 'ssid', title: 'SSID' },
    { key: 'bssid', title: 'BSSID' },
    {
      key: 'branch_name',
      title: 'Chi nhánh',
      render: (item: Wifi) => item.branch_name || item.branch?.name || '—',
    },
    {
      key: 'actions',
      title: 'Thao tác',
      width: '80px',
      render: (item: Wifi) => (
        <div className="cell-actions">
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
