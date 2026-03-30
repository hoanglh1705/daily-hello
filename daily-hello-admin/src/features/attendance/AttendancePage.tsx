import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getAttendances } from './api'
import type { Attendance } from './types'
import AttendanceTable from './components/AttendanceTable'

export default function AttendancePage() {
  const [data, setData] = useState<Attendance[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)
  const [date, setDate] = useState('')

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getAttendances({
        page,
        limit: DEFAULT_LIMIT,
        date: date || undefined,
      })
      setData(res.data)
      setTotal(res.meta.total)
    } catch (err) {
      console.error('Failed to fetch attendances', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, date])

  return (
    <div>
      <h1>Cham cong</h1>

      <div className="toolbar">
        <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
      </div>

      <AttendanceTable data={data} loading={loading} />

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />
    </div>
  )
}
