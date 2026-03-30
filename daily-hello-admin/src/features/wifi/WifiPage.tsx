import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import Modal from '@/shared/components/Modal'
import SearchSelect from '@/shared/components/SearchSelect'
import type { SearchSelectOption } from '@/shared/components/SearchSelect'
import { useDebounce } from '@/shared/hooks/useDebounce'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getWifiList, createWifi, deleteWifi } from './api'
import { getBranches } from '@/features/branch/api'
import type { Wifi } from './types'
import WifiTable from './components/WifiTable'
import WifiForm from './components/WifiForm'

export default function WifiPage() {
  const [data, setData] = useState<Wifi[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)

  const [branchId, setBranchId] = useState<number | null>(null)
  const [branchSearch, setBranchSearch] = useState('')
  const debouncedBranchSearch = useDebounce(branchSearch)
  const [branchOptions, setBranchOptions] = useState<SearchSelectOption[]>([])
  const [branchLoading, setBranchLoading] = useState(false)

  const [modalOpen, setModalOpen] = useState(false)

  const fetchBranches = async () => {
    setBranchLoading(true)
    try {
      const res = await getBranches({ page: 1, limit: 50, search: debouncedBranchSearch })
      setBranchOptions(res.data.items.map((b) => ({ value: b.id, label: b.name })))
    } catch (err) {
      console.error('Failed to fetch branches', err)
    } finally {
      setBranchLoading(false)
    }
  }

  useEffect(() => {
    fetchBranches()
  }, [debouncedBranchSearch])

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getWifiList({
        page,
        limit: DEFAULT_LIMIT,
        branch_id: branchId ?? undefined,
      })
      setData(res.data.items)
      setTotal(res.data.meta.total)
    } catch (err) {
      console.error('Failed to fetch wifi list', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, branchId])

  const handleBranchChange = (value: number | null) => {
    setBranchId(value)
    setPage(DEFAULT_PAGE)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Xac nhan xoa?')) return
    await deleteWifi(id)
    fetchData()
  }

  const handleSubmit = async (formData: {name: string, code: string, ssid: string; bssid: string; branch_id: number }) => {
    await createWifi(formData)
    setModalOpen(false)
    fetchData()
  }

  return (
    <div>
      <h1>Quan ly WiFi</h1>

      <div className="toolbar">
        <SearchSelect
          options={branchOptions}
          value={branchId}
          onChange={handleBranchChange}
          placeholder="Loc theo chi nhanh"
          searchValue={branchSearch}
          onSearchChange={setBranchSearch}
          loading={branchLoading}
        />
        <button onClick={() => setModalOpen(true)}>Them WiFi</button>
      </div>

      <WifiTable data={data} loading={loading} onDelete={handleDelete} />

      <Pagination
        page={page}
        limit={DEFAULT_LIMIT}
        total={total}
        onPageChange={setPage}
      />

      <Modal
        open={modalOpen}
        title="Them WiFi"
        onClose={() => setModalOpen(false)}
      >
        <WifiForm
          onSubmit={handleSubmit}
          onCancel={() => setModalOpen(false)}
        />
      </Modal>
    </div>
  )
}
