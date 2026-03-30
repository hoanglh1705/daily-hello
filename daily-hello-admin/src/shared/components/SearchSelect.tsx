import { useEffect, useRef, useState } from 'react'

export type SearchSelectOption = {
  value: number
  label: string
}

type Props = {
  options: SearchSelectOption[]
  value: number | null
  onChange: (value: number | null) => void
  placeholder?: string
  searchValue: string
  onSearchChange: (value: string) => void
  loading?: boolean
}

export default function SearchSelect({
  options,
  value,
  onChange,
  placeholder = 'Chon...',
  searchValue,
  onSearchChange,
  loading,
}: Props) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  const selectedLabel = options.find((o) => o.value === value)?.label

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  return (
    <div className="search-select" ref={ref}>
      <div className="search-select-trigger" onClick={() => setOpen(!open)}>
        <span className={selectedLabel ? '' : 'search-select-placeholder'}>
          {selectedLabel || placeholder}
        </span>
        {value != null && (
          <button
            className="search-select-clear"
            onClick={(e) => {
              e.stopPropagation()
              onChange(null)
              onSearchChange('')
            }}
          >
            &times;
          </button>
        )}
      </div>

      {open && (
        <div className="search-select-dropdown">
          <input
            className="search-select-input"
            placeholder="Tim kiem..."
            value={searchValue}
            onChange={(e) => onSearchChange(e.target.value)}
            autoFocus
          />
          <ul className="search-select-list">
            {loading && <li className="search-select-empty">Dang tai...</li>}
            {!loading && options.length === 0 && (
              <li className="search-select-empty">Khong co ket qua</li>
            )}
            {!loading &&
              options.map((opt) => (
                <li
                  key={opt.value}
                  className={`search-select-item ${opt.value === value ? 'selected' : ''}`}
                  onClick={() => {
                    onChange(opt.value)
                    setOpen(false)
                    onSearchChange('')
                  }}
                >
                  {opt.label}
                </li>
              ))}
          </ul>
        </div>
      )}
    </div>
  )
}
