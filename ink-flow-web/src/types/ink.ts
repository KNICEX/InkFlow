import { emptyUser, type User } from '@/types/user.ts'

export enum InkStatus {
  UnKnown = 0,
  UnPublished = 1,
  Pending = 2,
  Rejected = 3,
  Published = 4,
  Private = 5,
}

export const inkStatusLabel = (status: InkStatus) => {
  switch (status) {
    case InkStatus.UnPublished:
      return '未发布'
    case InkStatus.Pending:
      return '审核中'
    case InkStatus.Published:
      return '已发布'
    case InkStatus.Rejected:
      return '审核未通过'
    case InkStatus.Private:
      return '私有'
    default:
      return 'Unknown'
  }
}

export const inkStatusProp = (status: InkStatus) => {
  switch (status) {
    case InkStatus.UnPublished:
      return 'unPublished'
    case InkStatus.Pending:
      return 'pending'
    case InkStatus.Published:
      return 'published'
    case InkStatus.Rejected:
      return 'rejected'
    case InkStatus.Private:
      return 'private'
    default:
      return 'published'
  }
}

export const inkStatusFromProp = (prop: string) => {
  switch (prop) {
    case 'unPublished':
      return InkStatus.UnPublished
    case 'pending':
      return InkStatus.Pending
    case 'published':
      return InkStatus.Published
    case 'rejected':
      return InkStatus.Rejected
    case 'private':
      return InkStatus.Private
    default:
      return InkStatus.UnKnown
  }
}

export interface Category {
  id: number
  name: string
}

export const emptyCategory = (): Category => {
  return {
    id: 0,
    name: '',
  }
}

export interface Interactive {
  viewCnt: number
  likeCnt: number
  favoriteCnt: number
  commentCnt: number
  liked: boolean
  favorited: boolean
}

export const emptyInteractive = (): Interactive => {
  return {
    viewCnt: 0,
    likeCnt: 0,
    favoriteCnt: 0,
    commentCnt: 0,
    liked: false,
    favorited: false,
  }
}

export interface Ink {
  id: string
  cover: string
  title: string
  summary: string
  category: Category
  contentType: number
  contentHtml: string
  contentMeta: string

  status: InkStatus
  tags: string[]

  author: User
  interactive: Interactive

  createdAt: Date
  updatedAt: Date
}

export const emptyInk = (): Ink => {
  return {
    id: '0',
    cover: '',
    title: '',
    summary: '',
    category: emptyCategory(),
    contentType: 0,
    contentHtml: '',
    contentMeta: '',

    status: InkStatus.UnKnown,
    tags: [],
    author: emptyUser(),
    interactive: emptyInteractive(),

    // 0
    createdAt: new Date(0),
    updatedAt: new Date(0),
  }
}

export const getStatusLabel = (status: InkStatus) => {
  switch (status) {
    case InkStatus.UnPublished:
      return 'Draft'
    case InkStatus.Pending:
      return 'Pending'
    case InkStatus.Published:
      return 'Published'
    case InkStatus.Private:
      return 'Hidden'
  }
}
