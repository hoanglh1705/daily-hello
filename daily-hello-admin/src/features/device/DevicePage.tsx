import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import { DEFAULT_LIMIT, DEFAULT_PAGE } from '@/shared/utils/constants'
import { approveDevice, getDevices, rejectDevice } from './api'
import DeviceTable from './components/DeviceTable'
import type { Device, DeviceStatus } from './types'
import { getBranches } from '@/features/branch/api'
import type { Branch } from '@/features/branch/types'
import { getCurrentBranchId } from '@/services/tokenStorage'

const statusOptions: Array<{ value: DeviceStatus; label: string }> = [
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
]

export default function DevicePage() {
  const [data, setData] = useState<Device[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)
  const [status, setStatus] = useState<DeviceStatus>('pending')
  const [actionLoading, setActionLoading] = useState<number | null>(null)
  const [branchId, setBranchId] = useState<number | ''>(getCurrentBranchId() ?? '')
  const [branches, setBranches] = useState<Branch[]>([])

  const fetchDevices = async () => {
    setLoading(true)
    try {
      const res = await getDevices({
        page,
        limit: DEFAULT_LIMIT,
        status,
        branch_id: branchId || undefined,
      })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch devices', err)
    } finally {
      setLoading(false)
    }
  }

  const handleAction = async (id: number, action: () => Promise<unknown>) => {
    setActionLoading(id)
    try {
      await action()
      await fetchDevices()
    } catch (err) {
      console.error('Failed to update device status', err)
    } finally {
      setActionLoading(null)
    }
  }

  useEffect(() => {
    fetchDevices()
  }, [page, status, branchId])

  useEffect(() => {
    setPage(DEFAULT_PAGE)
  }, [status, branchId])

  useEffect(() => {
    getBranches({ page: 1, limit: 1000 })
      .then((res) => setBranches(res.data.items))
      .catch((err) => console.error('Failed to fetch branches', err))
  }, [])

  return (
    <div>
      <div className="page-header">
        <h1>Quan ly thiet bi</h1>
      </div>

      <div className="toolbar">
        <div className="toolbar-filters">
          <select
            value={branchId}
            onChange={(e) => setBranchId(e.target.value ? Number(e.target.value) : '')}
          >
            <option value="">Tat ca chi nhanh</option>
            {branches.map((branch) => (
              <option key={branch.id} value={branch.id}>
                {branch.name}
              </option>
            ))}
          </select>
          {statusOptions.map((option) => (
            <button
              key={option.value}
              className={`toolbar-chip ${status === option.value ? 'active' : ''}`}
              onClick={() => setStatus(option.value)}
            >
              {option.label}
            </button>
          ))}
        </div>
      </div>

      <DeviceTable
        data={data}
        loading={loading}
        actionLoading={actionLoading}
        onApprove={(device) => handleAction(device.id, () => approveDevice(device.id))}
        onReject={(device) => handleAction(device.id, () => rejectDevice(device.id))}
      />

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />
    </div>
  )
}
