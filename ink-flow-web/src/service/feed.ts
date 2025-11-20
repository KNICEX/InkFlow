import { get } from './axios.ts'
import type { Ink } from '@/types/ink.ts'

export const recommendInks = async (offset: number, limit: number) => {
  return await get<Ink[]>('/feed/ink/recommend', {
    params: {
      offset,
      limit,
    },
  })
}

export const followInks = async (maxId: string, timestamp: number, limit: number) => {
  return await get<Ink[]>('/feed/ink/follow', {
    params: {
      maxId,
      timestamp,
      limit,
    },
  })
}

export const hotInks = async (offset: number, limit: number) => {
  return await get<Ink[]>('/feed/ink/hot', {
    params: {
      offset,
      limit,
    },
  })
}
