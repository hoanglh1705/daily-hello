import axios from '@/services/axios'
import type { UserFilter } from './types'

export const getUsers = (params: UserFilter) => {
  return axios.get('/v1/users', { params })
}

export const getBranches = (params?: { page: number; limit: number; search?: string }) => {
  return axios.get('/v1/branches', { params })
}

export const createUser = (data: any) => {
  return axios.post('/v1/users', data)
}

export const updateUser = (id: number, data: any) => {
  return axios.put(`/v1/users/${id}`, data)
}
