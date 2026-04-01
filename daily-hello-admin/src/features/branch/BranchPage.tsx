import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import Modal from '@/shared/components/Modal'
import { useDebounce } from '@/shared/hooks/useDebounce'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getBranches, createBranch, updateBranch, deleteBranch } from './api'
import type { Branch } from './types'
import BranchTable from './components/BranchTable'
import BranchForm from './components/BranchForm'

export default function BranchPage() {
  const [data, setData] = useState<Branch[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search)

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Branch | null>(null)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getBranches({ page, limit: DEFAULT_LIMIT, search: debouncedSearch })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch branches', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, debouncedSearch])

  const handleCreate = () => {
    setEditing(null)
    setModalOpen(true)
  }

  const handleEdit = (branch: Branch) => {
    setEditing(branch)
    setModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Xác nhận xóa?')) return
    await deleteBranch(id)
    fetchData()
  }

  const handleSubmit = async (formData: { name: string; address: string; lat: number; lng: number; radius: number }) => {
    if (editing) {
      await updateBranch(editing.id, formData)
    } else {
      await createBranch(formData)
    }
    setModalOpen(false)
    fetchData()
  }

  return (
    <div>
      <div className="page-header">
        <h1>Quản lý chi nhánh</h1>
        <button className="btn-primary" onClick={handleCreate}>+ Thêm chi nhánh</button>
      </div>

      <div className="toolbar">
        <div className="toolbar-search">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
          <input
            placeholder="Tìm kiếm chi nhánh..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <BranchTable
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
        title={editing ? 'Cập nhật chi nhánh' : 'Thêm chi nhánh'}
        onClose={() => setModalOpen(false)}
      >
        <BranchForm
          initial={editing}
          onSubmit={handleSubmit}
          onCancel={() => setModalOpen(false)}
        />
      </Modal>
    </div>
  )
}
