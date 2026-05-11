import { useEffect, useState } from 'react'
import type { Wifi } from '../types'

type Props = {
  initial?: Wifi | null
  onSubmit: (data: { name: string; code: string; ssid: string; bssid: string; branch_id: number }) => void
  onCancel: () => void
}

export default function WifiForm({ initial, onSubmit, onCancel }: Props) {
  const [name, setName] = useState('')
  const [code, setCode] = useState('')
  const [ssid, setSsid] = useState('')
  const [bssid, setBssid] = useState('')
  const [branchId, setBranchId] = useState('')

  useEffect(() => {
    if (initial) {
      setName(initial.name)
      setCode(initial.code)
      setSsid(initial.ssid)
      setBssid(initial.bssid)
      setBranchId(String(initial.branch_id))
    } else {
      setName('')
      setCode('')
      setSsid('')
      setBssid('')
      setBranchId('')
    }
  }, [initial])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ name, code, ssid, bssid, branch_id: Number(branchId) })
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Name</label>
        <input value={name} onChange={(e) => setName(e.target.value)} required />
      </div>
      <div>
        <label>Code</label>
        <input value={code} onChange={(e) => setCode(e.target.value)} required />
      </div>
      <div>
        <label>SSID</label>
        <input value={ssid} onChange={(e) => setSsid(e.target.value)} required />
      </div>
      <div>
        <label>BSSID</label>
        <input value={bssid} onChange={(e) => setBssid(e.target.value)} required />
      </div>
      <div>
        <label>Branch ID</label>
        <input
          type="number"
          value={branchId}
          onChange={(e) => setBranchId(e.target.value)}
          required
        />
      </div>
      <div>
        <button type="submit">{initial ? 'Cập nhật' : 'Thêm mới'}</button>
        <button type="button" onClick={onCancel}>
          Hủy
        </button>
      </div>
    </form>
  )
}
