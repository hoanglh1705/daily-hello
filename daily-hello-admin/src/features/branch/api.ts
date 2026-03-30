import axios from '@/services/axios'
import type { Branch } from './types'

type ListResponse = {
  data: {
    items: Branch[]
    meta: { page: number; limit: number; total: number }
  }
}

export const getBranches = (params: { page: number; limit: number; search?: string }) => {
  return axios.get<unknown, ListResponse>('/v1/branches', { params })
}

export const createBranch = (data: Partial<Omit<Branch, 'id' | 'created_at'>>) => {
  return axios.post('/v1/branches', data)
}

export const updateBranch = (id: number, data: Partial<Omit<Branch, 'id' | 'created_at'>>) => {
  return axios.put(`/v1/branches/${id}`, data)
}

export const deleteBranch = (id: number) => {
  return axios.delete(`/v1/branches/${id}`)
}
