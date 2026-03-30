import { useEffect, useState } from 'react'
import Pagination from '@/shared/components/Pagination'
import Modal from '@/shared/components/Modal'
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '@/shared/utils/constants'
import { getWifiList, createWifi, deleteWifi } from './api'
import type { Wifi } from './types'
import WifiTable from './components/WifiTable'
import WifiForm from './components/WifiForm'

export default function WifiPage() {
  const [data, setData] = useState<Wifi[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(DEFAULT_PAGE)
  const [total, setTotal] = useState(0)

  const [modalOpen, setModalOpen] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getWifiList({ page, limit: DEFAULT_LIMIT })
      setData(res.data)
      setTotal(res.meta.total)
    } catch (err) {
      console.error('Failed to fetch wifi list', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page])

  const handleDelete = async (id: number) => {
    if (!confirm('Xac nhan xoa?')) return
    await deleteWifi(id)
    fetchData()
  }

  const handleSubmit = async (formData: { ssid: string; bssid: string; branch_id: number }) => {
    await createWifi(formData)
    setModalOpen(false)
    fetchData()
  }

  return (
    <div>
      <h1>Quan ly WiFi</h1>

      <div className="toolbar">
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
