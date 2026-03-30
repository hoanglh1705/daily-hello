import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import SearchSelect, { type SearchSelectOption } from '@/shared/components/SearchSelect'
import Modal from '@/shared/components/Modal'
import { useDebounce } from '@/shared/hooks/useDebounce'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getUsers, getBranches, createUser, updateUser } from './api'
import type { User } from './types'
import UserTable from './components/UserTable'
import UserForm from './components/UserForm'

export default function UserPage() {
  const [data, setData] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)

  // Filters
  const [keyword, setKeyword] = useState('')
  const debouncedKeyword = useDebounce(keyword)

  const [branchId, setBranchId] = useState<number | null>(null)

  // Branch dropdown state
  const [branchOptions, setBranchOptions] = useState<SearchSelectOption[]>([])
  const [branchSearch, setBranchSearch] = useState('')
  const debouncedBranchSearch = useDebounce(branchSearch)
  const [loadingBranches, setLoadingBranches] = useState(false)

  // Fetch users
  const fetchUsers = async () => {
    setLoading(true)
    try {
      const res = await getUsers({
        page,
        limit: DEFAULT_LIMIT,
        keyword: debouncedKeyword,
        branch_id: branchId,
      })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch users', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditing(null)
    setModalOpen(true)
  }

  const handleEdit = (user: User) => {
    setEditing(user)
    setModalOpen(true)
  }

  const handleSubmit = async (formData: any) => {
    if (editing) {
      await updateUser(editing.id, formData)
    } else {
      await createUser(formData)
    }
    setModalOpen(false)
    fetchUsers()
  }

  // Fetch branches for filter
  const fetchBranches = async () => {
    setLoadingBranches(true)
    try {
      // Fetch matching branches to show in dropdown
      const res = await getBranches({ page: 1, limit: 100, search: debouncedBranchSearch })
      const opts = res.data.items.map((b: any) => ({ value: b.id, label: b.name }))
      setBranchOptions(opts)
    } catch (err) {
      console.error('Failed to fetch branches', err)
    } finally {
      setLoadingBranches(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [page, debouncedKeyword, branchId])

  useEffect(() => {
    fetchBranches()
  }, [debouncedBranchSearch])

  // Reset page when filter changes
  useEffect(() => {
    setPage(DEFAULT_PAGE)
  }, [debouncedKeyword, branchId])

  return (
    <div>
      <h1>Quan ly user</h1>

      <div className="toolbar" style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
        <input
          placeholder="Tim kiem keyword..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          style={{ width: '250px' }}
        />

        <div style={{ width: '250px' }}>
          <SearchSelect
            options={branchOptions}
            value={branchId}
            onChange={setBranchId}
            placeholder="Loc theo chi nhanh..."
            searchValue={branchSearch}
            onSearchChange={setBranchSearch}
            loading={loadingBranches}
          />
        </div>

        <button onClick={handleCreate}>Thêm User</button>
      </div>

      <UserTable data={data} loading={loading} onEdit={handleEdit} />

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />

      <Modal
        open={modalOpen}
        title={editing ? 'Chỉnh sửa User' : 'Thêm mới User'}
        onClose={() => setModalOpen(false)}
      >
        <UserForm
          initial={editing}
          onSubmit={handleSubmit}
          onCancel={() => setModalOpen(false)}
        />
      </Modal>
    </div>
  )
}
