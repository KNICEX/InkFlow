import { get } from './axios.ts'
import type { TagStats } from '../types/stats.ts'
export const topTag = (offset: number, limit: number) => {
  return get<TagStats[]>('/stats/top-tags', {
    params: {
      offset,
      limit,
    },
  })
}
