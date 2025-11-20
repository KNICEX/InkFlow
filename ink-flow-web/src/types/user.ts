import type { Ink } from '@/types/ink.ts'

export interface User {
  id: string
  avatar: string
  banner: string
  username: string
  account: string
  aboutMe: string
  email: string
  phone: string
  links: string[]

  followed: boolean
  followers: number
  following: number

  createdAt: Date
  updatedAt: Date
}

export interface UserFollow extends User {
  inks: Ink[]
}

export interface StatisticInfo {
  inkCount: number
  commentCount: number
  followerCount: number
  followingCount: number
  likedCount: number
  likedInkCount: number
  collectionCount: number
  readCount: number
  joinedDays: number
}

export const emptyUser = (): User => {
  return {
    id: '0',
    avatar: '',
    banner: '',
    username: '',
    account: '',
    aboutMe: '',
    email: '',
    phone: '',
    links: [],
    followed: false,
    followers: 0,
    following: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
  }
}
