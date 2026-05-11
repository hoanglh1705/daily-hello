import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import Modal from '@/shared/components/Modal'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getHolidays, createHoliday, updateHoliday, deleteHoliday } from './api'
import type { Holiday } from './types'
import HolidayTable from './components/HolidayTable'
import HolidayForm from './components/HolidayForm'

export default function HolidayPage() {
  const [data, setData] = useState<Holiday[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Holiday | null>(null)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getHolidays({ page, limit: DEFAULT_LIMIT })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch holidays', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page])

  const handleCreate = () => {
    setEditing(null)
    setModalOpen(true)
  }

  const handleEdit = (holiday: Holiday) => {
    setEditing(holiday)
    setModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Xác nhận xóa ngày nghỉ này?')) return
    await deleteHoliday(id)
    fetchData()
  }

  const handleSubmit = async (formData: { name: string; date: string; description: string }) => {
    if (editing) {
      await updateHoliday(editing.id, formData)
    } else {
      await createHoliday(formData)
    }
    setModalOpen(false)
    fetchData()
  }

  return (
    <div>
      <div className="page-header">
        <h1>Quản lý ngày nghỉ</h1>
        <button className="btn-primary" onClick={handleCreate}>+ Thêm ngày nghỉ</button>
      </div>

      <HolidayTable
        data={data}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />

      <Modal
        open={modalOpen}
        title={editing ? 'Cập nhật ngày nghỉ' : 'Thêm ngày nghỉ'}
        onClose={() => setModalOpen(false)}
      >
        <HolidayForm
          initial={editing}
          onSubmit={handleSubmit}
          onCancel={() => setModalOpen(false)}
        />
      </Modal>
    </div>
  )
}
