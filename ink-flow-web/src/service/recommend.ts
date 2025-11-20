import type { User } from '@/types/user.ts'
import { get } from './axios.ts'
import type { Ink } from '@/types/ink.ts'

export const recommendAuthor = async (offset: number, limit: number) => {
  return await get<User[]>('/recommend/author', {
    params: {
      offset,
      limit,
    },
  })
}

export const similarInks = async (id: string, offset: number, limit: number) => {
  return await get<Ink[]>(`/recommend/ink/similar/${id}`, {
    params: {
      offset,
      limit,
    },
  })
}

export const similarAuthors = async (id: string, offset: number, limit: number) => {
  return await get<User[]>(`/recommend/author/similar/${id}`, {
    params: {
      offset,
      limit,
    },
  })
}
