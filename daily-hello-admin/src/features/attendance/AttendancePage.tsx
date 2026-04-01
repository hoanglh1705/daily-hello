import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import {
  approveCheckIn,
  approveCheckOut,
  getAttendances,
  rejectCheckIn,
  rejectCheckOut,
} from './api'
import type { Attendance } from './types'
import AttendanceTable from './components/AttendanceTable'
import { getBranches } from '../branch/api'
import type { Branch } from '../branch/types'
import { getCurrentBranchId } from '@/services/tokenStorage'

export default function AttendancePage() {
  const [data, setData] = useState<Attendance[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)
  const [branchId, setBranchId] = useState<number | ''>(getCurrentBranchId() ?? '')
  const [fromDate, setFromDate] = useState('')
  const [toDate, setToDate] = useState('')
  const [branches, setBranches] = useState<Branch[]>([])
  const [actionLoading, setActionLoading] = useState<number | null>(null)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getAttendances({
        page,
        limit: DEFAULT_LIMIT,
        branch_id: branchId || undefined,
        from: fromDate || undefined,
        to: toDate || undefined,
      })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch attendances', err)
    } finally {
      setLoading(false)
    }
  }

  const fetchBranches = async () => {
    try {
      const res = await getBranches({ page: 1, limit: 1000 })
      setBranches(res.data.items)
    } catch (err) {
      console.error('Failed to fetch branches', err)
    }
  }

  const handleAttendanceAction = async (id: number, action: () => Promise<unknown>) => {
    setActionLoading(id)
    try {
      await action()
      await fetchData()
    } catch (err) {
      console.error('Failed to update attendance status', err)
    } finally {
      setActionLoading(null)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, branchId, fromDate, toDate])

  useEffect(() => {
    fetchBranches()
  }, [])

  return (
    <div>
      <h1>Cham cong</h1>

      <div className="toolbar">
        <select
          value={branchId}
          onChange={(e) => {
            setBranchId(e.target.value ? Number(e.target.value) : '')
            setPage(DEFAULT_PAGE)
          }}
        >
          <option value="">Tat ca chi nhanh</option>
          {branches.map((branch) => (
            <option key={branch.id} value={branch.id}>
              {branch.name}
            </option>
          ))}
        </select>
        <input
          type="date"
          value={fromDate}
          onChange={(e) => {
            setFromDate(e.target.value)
            setPage(DEFAULT_PAGE)
          }}
        />
        <input
          type="date"
          value={toDate}
          onChange={(e) => {
            setToDate(e.target.value)
            setPage(DEFAULT_PAGE)
          }}
        />
      </div>

      <AttendanceTable
        data={data}
        loading={loading}
        actionLoading={actionLoading}
        onApproveCheckIn={(item) =>
          handleAttendanceAction(item.id, () => approveCheckIn(item.id))
        }
        onRejectCheckIn={(item) =>
          handleAttendanceAction(item.id, () => rejectCheckIn(item.id))
        }
        onApproveCheckOut={(item) =>
          handleAttendanceAction(item.id, () => approveCheckOut(item.id))
        }
        onRejectCheckOut={(item) =>
          handleAttendanceAction(item.id, () => rejectCheckOut(item.id))
        }
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
