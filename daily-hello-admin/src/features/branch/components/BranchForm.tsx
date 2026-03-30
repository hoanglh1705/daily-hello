import { useState, useEffect } from 'react'
import type { Branch } from '../types'

type Props = {
  initial?: Branch | null
  onSubmit: (data: { name: string; address: string }) => void
  onCancel: () => void
}

export default function BranchForm({ initial, onSubmit, onCancel }: Props) {
  const [name, setName] = useState('')
  const [address, setAddress] = useState('')

  useEffect(() => {
    if (initial) {
      setName(initial.name)
      setAddress(initial.address)
    } else {
      setName('')
      setAddress('')
    }
  }, [initial])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ name, address })
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Ten chi nhanh</label>
        <input value={name} onChange={(e) => setName(e.target.value)} required />
      </div>
      <div>
        <label>Dia chi</label>
        <input value={address} onChange={(e) => setAddress(e.target.value)} required />
      </div>
      <div>
        <button type="submit">{initial ? 'Cap nhat' : 'Tao moi'}</button>
        <button type="button" onClick={onCancel}>
          Huy
        </button>
      </div>
    </form>
  )
}
