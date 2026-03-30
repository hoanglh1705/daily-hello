export type PaginatedListResponse<T> = {
  success?: boolean
  data: {
    items: T[]
    meta: { page: number; limit: number; total: number }
  }
  error_code?: string
  error_message?: string
}
