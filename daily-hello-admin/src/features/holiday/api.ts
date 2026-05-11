import axios from '@/services/axios'
import type { Holiday } from './types'

type ListResponse = {
  data: {
    items: Holiday[]
    meta: { page: number; limit: number; total: number }
  }
}

type HolidayListResponse = {
  data: Holiday[]
}

export const getHolidays = (params: { page: number; limit: number }) => {
  return axios.get<unknown, ListResponse>('/v1/admin/holidays', { params })
}

export const createHoliday = (data: { name: string; date: string; description?: string }) => {
  return axios.post('/v1/admin/holidays', data)
}

export const updateHoliday = (id: number, data: { name?: string; date?: string; description?: string }) => {
  return axios.put(`/v1/admin/holidays/${id}`, data)
}

export const deleteHoliday = (id: number) => {
  return axios.delete(`/v1/admin/holidays/${id}`)
}

export const getHolidaysByMonth = (params: { year: number; month: number }) => {
  return axios.get<unknown, HolidayListResponse>('/v1/holidays', { params })
}
