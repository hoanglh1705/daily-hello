import Table from '@/shared/components/Table'
import { formatDateTime } from '@/shared/utils/formatDate'
import type { Device } from '../types'

type Props = {
  data: Device[]
  loading: boolean
  actionLoading?: number | null
  onApprove?: (device: Device) => void
  onReject?: (device: Device) => void
}

function getInitials(name: string) {
  return name
    .split(' ')
    .map((word) => word[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

function StatusBadge({ status }: { status: string }) {
  const classMap: Record<string, string> = {
    pending: 'pending',
    approved: 'active',
    rejected: 'on-leave',
  }

  return (
    <span className={`status-badge ${classMap[status] || 'inactive'}`}>
      {status}
    </span>
  )
}

export default function DeviceTable({
  data,
  loading,
  actionLoading,
  onApprove,
  onReject,
}: Props) {
  const columns = [
    { key: 'id', title: 'ID', width: '72px' },
    {
      key: 'user',
      title: 'Nguoi dung',
      render: (device: Device) => {
        const user = device.user
        if (!user) {
          return '—'
        }

        return (
          <div className="cell-user">
            <div className="cell-avatar">{getInitials(user.name)}</div>
            <div className="cell-user-info">
              <span className="cell-user-name">{user.name}</span>
              <span className="cell-user-sub">{user.email}</span>
            </div>
          </div>
        )
      },
    },
    { key: 'device_name', title: 'Ten thiet bi' },
    { key: 'platform', title: 'Nen tang' },
    { key: 'model', title: 'Model' },
    {
      key: 'device_id',
      title: 'Device ID',
      render: (device: Device) => (
        <code className="device-code">{device.device_id || '—'}</code>
      ),
    },
    {
      key: 'status',
      title: 'Trang thai',
      render: (device: Device) => <StatusBadge status={device.status} />,
    },
    {
      key: 'created_at',
      title: 'Dang ky luc',
      render: (device: Device) => formatDateTime(device.created_at),
    },
    {
      key: 'approved_at',
      title: 'Xu ly luc',
      render: (device: Device) =>
        device.approved_at ? formatDateTime(device.approved_at) : '—',
    },
    {
      key: 'action',
      title: 'Thao tac',
      width: '160px',
      render: (device: Device) => {
        const isPending = device.status === 'pending'
        const isActing = actionLoading === device.id

        return (
          <div className="cell-actions">
            <button
              className="device-action-btn device-action-btn-approve"
              onClick={() => onApprove?.(device)}
              disabled={!isPending || isActing}
            >
              Approve
            </button>
            <button
              className="device-action-btn device-action-btn-reject"
              onClick={() => onReject?.(device)}
              disabled={!isPending || isActing}
            >
              Reject
            </button>
          </div>
        )
      },
    },
  ]

  return <Table data={data} columns={columns} loading={loading} />
}
