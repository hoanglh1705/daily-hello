import { useState } from 'react'

type Props = {
  onSubmit: (data: { ssid: string; bssid: string; branch_id: number }) => void
  onCancel: () => void
}

export default function WifiForm({ onSubmit, onCancel }: Props) {
  const [ssid, setSsid] = useState('')
  const [bssid, setBssid] = useState('')
  const [branchId, setBranchId] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ ssid, bssid, branch_id: Number(branchId) })
  }

  return (
    <form onSubmit={handleSubmit}>
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
        <button type="submit">Tao moi</button>
        <button type="button" onClick={onCancel}>
          Huy
        </button>
      </div>
    </form>
  )
}
