import Table from '@/shared/components/Table'
import { formatDateTime } from '@/shared/utils/formatDate'
import type { Attendance } from '../types'

type Props = {
  data: Attendance[]
  loading: boolean
  actionLoading?: number | null
  onApproveCheckIn?: (item: Attendance) => void
  onRejectCheckIn?: (item: Attendance) => void
  onApproveCheckOut?: (item: Attendance) => void
  onRejectCheckOut?: (item: Attendance) => void
}

function renderLocation(lat: number | null, lng: number | null) {
  if (lat == null || lng == null) {
    return '-'
  }

  const coordinates = `${lat}, ${lng}`
  const mapUrl = `https://www.google.com/maps?q=${lat},${lng}`

  return (
    <a href={mapUrl} target="_blank" rel="noreferrer">
      {coordinates}
    </a>
  )
}

function renderAttendanceCard({
  tone,
  item,
  time,
  lat,
  lng,
  type,
  status,
  actionLoading,
  onApprove,
  onReject,
}: {
  tone: 'in' | 'out'
  item: Attendance
  time: string | null
  lat: number | null
  lng: number | null
  type: string | null
  status: string | null
  actionLoading?: number | null
  onApprove?: (item: Attendance) => void
  onReject?: (item: Attendance) => void
}) {
  const isPending = status === 'waiting_approve'
  const isActing = actionLoading === item.id

  return (
    <div className={`attendance-card attendance-card-${tone}`}>
      <div className="attendance-card-row">
        <span className="attendance-card-label">Thời gian</span>
        <strong>{time ? formatDateTime(time) : '-'}</strong>
      </div>
      <div className="attendance-card-row">
        <span className="attendance-card-label">Vị trí</span>
        {renderLocation(lat, lng)}
      </div>
      <div className="attendance-card-row">
        <span className="attendance-card-label">Loại</span>
        <span>{type || '-'}</span>
      </div>
      <div className="attendance-card-row">
        <span className="attendance-card-label">Trạng thái</span>
        <span>{status || '-'}</span>
      </div>
      {isPending && (
        <div className="attendance-card-actions">
          <button onClick={() => onApprove?.(item)} disabled={isActing}>
            Approve
          </button>
          <button
            className="attendance-card-reject"
            onClick={() => onReject?.(item)}
            disabled={isActing}
          >
            Reject
          </button>
        </div>
      )}
    </div>
  )
}

export default function AttendanceTable({
  data,
  loading,
  actionLoading,
  onApproveCheckIn,
  onRejectCheckIn,
  onApproveCheckOut,
  onRejectCheckOut,
}: Props) {
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
      key: 'check_in',
      title: 'Check in',
      render: (item: Attendance) =>
        renderAttendanceCard({
          tone: 'in',
          item,
          time: item.check_in_time,
          lat: item.check_in_lat,
          lng: item.check_in_lng,
          type: item.check_in_type,
          status: item.check_in_status,
          actionLoading,
          onApprove: onApproveCheckIn,
          onReject: onRejectCheckIn,
        }),
    },
    {
      key: 'check_out',
      title: 'Check out',
      render: (item: Attendance) =>
        renderAttendanceCard({
          tone: 'out',
          item,
          time: item.check_out_time,
          lat: item.check_out_lat,
          lng: item.check_out_lng,
          type: item.check_out_type,
          status: item.check_out_status,
          actionLoading,
          onApprove: onApproveCheckOut,
          onReject: onRejectCheckOut,
        }),
    },
  ]

  return (
    <Table
      columns={columns}
      data={data}
      loading={loading}
      getRowClassName={(_, index) =>
        index % 2 === 1 ? 'shared-table-row-even' : undefined
      }
    />
  )
}
