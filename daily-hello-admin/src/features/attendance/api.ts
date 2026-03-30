import axios from '@/services/axios'
import type { Attendance } from './types'
import type { PaginatedListResponse } from '@/shared/types/api'

type ListResponse = PaginatedListResponse<Attendance>

export const getAttendances = (params: {
  page: number
  limit: number
  branch_id?: number
  from?: string
  to?: string
}) => {
  return axios.get<unknown, ListResponse>('/v1/attendance/history', { params })
}

export const approveCheckIn = (id: number) => {
  return axios.put(`/v1/admin/attendance/${id}/check-in/approve`)
}

export const rejectCheckIn = (id: number) => {
  return axios.put(`/v1/admin/attendance/${id}/check-in/reject`)
}

export const approveCheckOut = (id: number) => {
  return axios.put(`/v1/admin/attendance/${id}/check-out/approve`)
}

export const rejectCheckOut = (id: number) => {
  return axios.put(`/v1/admin/attendance/${id}/check-out/reject`)
}
