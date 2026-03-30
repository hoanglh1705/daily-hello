import axios from '@/services/axios'
import type { Attendance } from './types'

type ListResponse = {
  data: Attendance[]
  meta: { page: number; limit: number; total: number }
}

export const getAttendances = (params: {
  page: number
  limit: number
  branch_id?: number
  date?: string
}) => {
  return axios.get<unknown, ListResponse>('/v1/attendances', { params })
}
