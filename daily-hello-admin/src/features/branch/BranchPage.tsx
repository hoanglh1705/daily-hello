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
    if (!confirm('Xac nhan xoa?')) return
    await deleteBranch(id)
    fetchData()
  }

  const handleSubmit = async (formData: { name: string; address: string }) => {
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
      <h1>Quan ly chi nhanh</h1>

      <div className="toolbar">
        <input
          placeholder="Tim kiem..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button onClick={handleCreate}>Them chi nhanh</button>
      </div>

      <BranchTable
        data={data}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />v1/branch-wifi

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />

      <Modal
        open={modalOpen}
        title={editing ? 'Cap nhat chi nhanh' : 'Them chi nhanh'}
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
