import { useState, useEffect } from 'react'
import type { Holiday } from '../types'

type Props = {
  initial: Holiday | null
  onSubmit: (data: { name: string; date: string; description: string }) => void
  onCancel: () => void
}

export default function HolidayForm({ initial, onSubmit, onCancel }: Props) {
  const [name, setName] = useState('')
  const [date, setDate] = useState('')
  const [description, setDescription] = useState('')

  useEffect(() => {
    if (initial) {
      setName(initial.name)
      // Convert date to YYYY-MM-DD format for input
      setDate(initial.date.split('T')[0])
      setDescription(initial.description || '')
    } else {
      setName('')
      setDate('')
      setDescription('')
    }
  }, [initial])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ name, date, description })
  }

  return (
    <form onSubmit={handleSubmit} className="form">
      <div className="form-group">
        <label htmlFor="holiday-name">Tên ngày nghỉ *</label>
        <input
          id="holiday-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="VD: Tết Nguyên Đán, Quốc khánh..."
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="holiday-date">Ngày *</label>
        <input
          id="holiday-date"
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="holiday-description">Mô tả</label>
        <textarea
          id="holiday-description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Mô tả thêm về ngày nghỉ..."
          rows={3}
        />
      </div>

      <div className="form-actions">
        <button type="button" className="btn-secondary" onClick={onCancel}>
          Hủy
        </button>
        <button type="submit" className="btn-primary">
          {initial ? 'Cập nhật' : 'Thêm mới'}
        </button>
      </div>
    </form>
  )
}
