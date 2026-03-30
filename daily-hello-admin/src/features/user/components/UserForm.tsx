import { useState, useEffect } from 'react'
import type { User } from '../types'
import SearchSelect, { type SearchSelectOption } from '@/shared/components/SearchSelect'
import { getBranches } from '../api'

type Props = {
  initial: User | null
  onSubmit: (data: any) => Promise<void>
  onCancel: () => void
}

export default function UserForm({ initial, onSubmit, onCancel }: Props) {
  const [loading, setLoading] = useState(false)
  
  // Basic states
  const [name, setName] = useState(initial?.name || '')
  const [code, setCode] = useState(initial?.code || '')
  const [email, setEmail] = useState(initial?.email || '')
  const [phone, setPhone] = useState(initial?.phone || '')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState(initial?.role || 'employee')
  const [status, setStatus] = useState(initial?.status || 'active')
  const [branchId, setBranchId] = useState<number | null>(initial?.branch_id || null)

  // Branch dropdown logic inside form
  const [branchOptions, setBranchOptions] = useState<SearchSelectOption[]>([])
  const [branchSearch, setBranchSearch] = useState('')
  const [loadingBranches, setLoadingBranches] = useState(false)

  useEffect(() => {
    const fetchBranches = async () => {
      setLoadingBranches(true)
      try {
        const res = await getBranches({ page: 1, limit: 100, search: branchSearch })
        const opts = res.data.items.map((b: any) => ({ value: b.id, label: b.name }))
        // Include initial branch if it's missing from search results and search is empty
        if (!branchSearch && initial?.branch_id && initial?.branch?.name) {
          if (!opts.find((o: any) => o.value === initial.branch_id)) {
            opts.unshift({ value: initial.branch_id, label: initial.branch.name })
          }
        }
        setBranchOptions(opts)
      } catch (err) {
        console.error('Failed to fetch branches', err)
      } finally {
        setLoadingBranches(false)
      }
    }
    
    // Add simple debounce manually just for this form or trigger on search
    const timer = setTimeout(fetchBranches, 300)
    return () => clearTimeout(timer)
  }, [branchSearch, initial])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const payload: any = {
        name,
        phone,
        role,
        branch_id: branchId || null,
      }
      if (initial) {
        // Edit mode (status valid on update)
        payload.status = status
      } else {
        // Create mode
        payload.code = code
        payload.email = email
        payload.password = password
      }
      await onSubmit(payload)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="form-container">
      <div className="form-group">
        <label>Tên user (*)</label>
        <input required value={name} onChange={(e) => setName(e.target.value)} />
      </div>

      {!initial && (
        <>
          <div className="form-group">
            <label>Mã user (*)</label>
            <input required value={code} onChange={(e) => setCode(e.target.value)} />
          </div>
          <div className="form-group">
            <label>Email (*)</label>
            <input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>
          <div className="form-group">
            <label>Mật khẩu (*)</label>
            <input type="password" required={!initial} value={password} onChange={(e) => setPassword(e.target.value)} />
          </div>
        </>
      )}

      <div className="form-group">
        <label>Điện thoại</label>
        <input value={phone} onChange={(e) => setPhone(e.target.value)} />
      </div>

      <div className="form-group">
        <label>Quyền (*)</label>
        <select value={role} onChange={(e) => setRole(e.target.value)}>
          <option value="admin">Admin</option>
          <option value="manager">Manager</option>
          <option value="employee">Employee</option>
        </select>
      </div>

      {initial && (
        <div className="form-group">
          <label>Trạng thái</label>
          <select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      )}

      <div className="form-group">
        <label>Chi nhánh</label>
        <SearchSelect
          options={branchOptions}
          value={branchId}
          onChange={setBranchId}
          placeholder="Chọn chi nhánh..."
          searchValue={branchSearch}
          onSearchChange={setBranchSearch}
          loading={loadingBranches}
        />
      </div>

      <div className="form-actions">
        <button type="button" onClick={onCancel} disabled={loading} className="btn-cancel">Hủy</button>
        <button type="submit" disabled={loading} className="btn-primary">Lưu</button>
      </div>
      
      {/* Some simple inline styles if missing from global, keeping it functional */}
      <style>{`
        .form-container { display: flex; flex-direction: column; gap: 12px; }
        .form-group { display: flex; flex-direction: column; gap: 4px; }
        .form-group input, .form-group select { padding: 8px; border: 1px solid #ccc; border-radius: 4px; }
        .form-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
        .btn-primary { background: #1677ff; color: #fff; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; }
        .btn-cancel { background: #f5f5f5; border: 1px solid #d9d9d9; padding: 8px 16px; border-radius: 4px; cursor: pointer; }
      `}</style>
    </form>
  )
}
