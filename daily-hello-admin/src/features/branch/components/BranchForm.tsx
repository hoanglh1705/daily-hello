import { useState, useEffect } from 'react'
import type { Branch } from '../types'

type Props = {
  initial?: Branch | null
  onSubmit: (data: { name: string; address: string; lat: number; lng: number; radius: number }) => void
  onCancel: () => void
}

export default function BranchForm({ initial, onSubmit, onCancel }: Props) {
  const [name, setName] = useState('')
  const [address, setAddress] = useState('')
  const [latLng, setLatLng] = useState('')
  const [latLngError, setLatLngError] = useState('')
  const [radius, setRadius] = useState('0')

  useEffect(() => {
    if (initial) {
      setName(initial.name)
      setAddress(initial.address)
      setLatLng(`${initial.lat}, ${initial.lng}`)
      setRadius(String(initial.radius))
    } else {
      setName('')
      setAddress('')
      setLatLng('')
      setRadius('0')
    }
  }, [initial])

  const parseLatLng = (value: string): [number, number] | null => {
    const parts = value.split(',')
    if (parts.length !== 2) return null
    const lat = parseFloat(parts[0].trim())
    const lng = parseFloat(parts[1].trim())
    if (isNaN(lat) || isNaN(lng)) return null
    return [lat, lng]
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const parsed = parseLatLng(latLng)
    if (!parsed) {
      setLatLngError('Dinh dang khong hop le. Vi du: 10.8029527674255, 106.61170066134497')
      return
    }
    setLatLngError('')
    onSubmit({ name, address, lat: parsed[0], lng: parsed[1], radius: parseFloat(radius) })
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
        <label>Lat, Lng</label>
        <input
          value={latLng}
          onChange={(e) => { setLatLng(e.target.value); setLatLngError('') }}
          placeholder="ex: 10.8029527674255, 106.61170066134497"
          required
        />
        {latLngError && <span style={{ color: 'red', fontSize: '0.85em' }}>{latLngError}</span>}
      </div>
      <div>
        <label>Radius(m)</label>
        <input value={radius} onChange={(e) => setRadius(e.target.value)} required />
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
