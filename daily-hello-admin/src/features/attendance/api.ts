import axios from '@/services/axios'
import type { Attendance } from './types'
import type { PaginatedListResponse } from '@/shared/types/api'

type ListResponse = PaginatedListResponse<Attendance>

export const getAttendances = (params: {
  page: number
  limit: number
  branch_id?: number
  date?: string
}) => {
  return axios.get<unknown, ListResponse>('/v1/attendance/history', { params })
}
