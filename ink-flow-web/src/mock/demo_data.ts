import type { Ink } from '@/types/ink.ts'
import { emptyInk, InkStatus } from '@/types/ink.ts'
import type { TagStats } from '@/types/stats.ts'
import type { User } from '@/types/user.ts'
import { emptyUser } from '@/types/user.ts'
import { type Comment, emptyComment } from '@/types/comment.ts'
import {
  emptyMergedLike,
  emptyNotification,
  NotificationType,
  SubjectType,
} from '@/types/notification.ts'

const randomDate = (days?: number) => {
  const now = new Date()
  return new Date(now.getTime() - Math.random() * (days || 365) * 24 * 60 * 60 * 1000)
}

const randomCover = (mustImage?: boolean) => {
  const urls = [
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742147867/qzr5v0bavnefmijyjhzd.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742200924/nchghsqliegrlcqmuqd1.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742156516/iedbssljjhobe5oltyrd.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742145347/j1ozvllqd0xe7dtw2x7c.png',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1741972579/twitter_snowmint_snowlistenz__20241003-030616_1841676112663114079_photo_rii427.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1741537291/twitter_%E8%B0%B7%E5%B7%9D_kibishiihiekomi__20250204-093016_1886708845923459466_photo_dntbvf.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1741704737/twitter_%E5%91%BD_atmarkaoi__20230511-085828_1656584545968484354_photo-3_mzi7ql.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742146011/pmbdmoo0dwf25wetphkl.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742148009/fdhqgjp7djhdyuzmps18.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742156721/cb9h1ixah2vgd3ytydax.jpg',
    'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1742150049/rxzafga7q1rx1y6cqdjp.jpg',
  ]
  if (!mustImage && Math.random() < 0.1) {
    return ''
  }
  return urls[Math.floor(Math.random() * urls.length)]
}

const randomTitle = () => {
  const titles = [
    'Hello this is a title',
    'This is a title',
    'This is a very very very long long long title',
    'Today is a good day',
    'What is web3',
    'What is blockchain',
    'What is NFT',
    'about the golang',
    'about the long long long long rust',
    'about the vue',
    'how to draw a cat',
    'how to draw a dog',
    'how to draw a bird',
  ]
  return titles[Math.floor(Math.random() * titles.length)]
}

const randomAccount = () => {
  const accounts = [
    'mikan',
    'cat',
    'misaki',
    'chtholly',
    'nephren',
    'willem',
    'ren',
    'kutori',
    'kazuki',
    'chise',
    'elias',
    'satoshi',
  ]
  return accounts[Math.floor(Math.random() * accounts.length)]
}

export const demoInteractive = () => {
  return {
    viewCnt: Math.floor(Math.random() * 1000),
    favoriteCnt: Math.floor(Math.random() * 1000),
    likeCnt: Math.floor(Math.random() * 1000),
    commentCnt: Math.floor(Math.random() * 1000),
    liked: Math.random() > 0.5,
    favorited: Math.random() > 0.5,
  }
}

export const demoInks = (n?: number) => {
  const res: Ink[] = []
  n = n == undefined ? 20 : n
  for (let i = 0; i < n; i++) {
    const ink = emptyInk()
    ink.id = i.toString()
    ink.cover = randomCover()
    ink.contentHtml = 'Hello this'
    ink.title = randomTitle()
    ink.tags = ['test', 'test2']
    ink.status = InkStatus.Published
    ink.author.avatar = randomCover()
    ink.author.username = 'Mikan'
    ink.author.account = randomAccount()
    ink.author.aboutMe = 'I am a cat'
    ink.author.followers = Math.floor(Math.random() * 1000)
    ink.author.following = Math.floor(Math.random() * 1000)
    ink.interactive = demoInteractive()
    ink.createdAt = randomDate(1)

    res.push(ink)
  }
  return res
}

export const demoUsers = (n?: number) => {
  const res: User[] = []
  n = n == undefined ? 20 : n

  const idStart = Math.floor(Math.random() * 100000)
  for (let i = 0; i < n; i++) {
    const user = emptyUser()
    user.id = (i + idStart).toString()
    user.avatar = randomCover()
    user.banner =
      'https://res.cloudinary.com/dw3cnfcb5/image/upload/v1741972579/twitter_snowmint_snowlistenz__20241003-030616_1841676112663114079_photo_rii427.jpg'
    user.username = randomAccount()
    user.account = randomAccount()
    user.aboutMe = 'I am a cat'
    user.email = ''
    user.phone = ''
    user.followed = Math.random() > 0.7
    user.followers = 100
    user.following = 100
    user.createdAt = randomDate(10)
    user.updatedAt = new Date()
    res.push(user)
  }
  return res
}

export const demoCommentPayload = () => {
  return {
    content: 'Hello this is a comment 😭',
    images: Math.random() > 0.7 ? [randomCover(true), randomCover(true)] : [],
  }
}

export const demoComments = (n?: number, sub?: boolean) => {
  const res: Comment[] = []
  const startId = Math.floor(Math.random() * 100000)
  n = n == undefined ? 20 : n
  for (let i = 0; i < n; i++) {
    const comment = emptyComment()
    comment.id = (i + startId).toString()
    comment.biz = 'ink'
    comment.bizId = '0'
    comment.commentator = demoUsers(1)[0]
    comment.isAuthor = Math.random() > 0.7
    comment.payload = demoCommentPayload()
    comment.parent = null
    comment.root = null
    comment.children = []
    comment.stats = {
      liked: Math.random() > 0.5,
      likeCnt: Math.floor(Math.random() * 1000),
      viewCnt: Math.floor(Math.random() * 1000),
      replyCnt: Math.floor(Math.random() * 1000),
    }
    comment.createdAt = new Date()
    res.push(comment)
  }

  if (sub) {
    for (const com of res) {
      if (Math.random() > 0.5) {
        com.children = demoComments(Math.floor(Math.random() * 5))
        for (const child of com.children) {
          child.parent = com
          child.root = com
        }
      }
    }
  }

  return res
}

export const demoNotifications = (type: NotificationType, n?: number) => {
  switch (type) {
    case NotificationType.Follow:
      const follows = []
      n = n == undefined ? 20 : n
      for (let i = 0; i < n; i++) {
        const no = emptyNotification(type)
        no.id = i.toString()
        no.subjectType = SubjectType.User
        no.user = demoUsers(1)[0]
        follows.push(no)
      }
      return follows
    case NotificationType.Like:
      const likes = []
      n = n == undefined ? 20 : n
      for (let i = 0; i < n; i++) {
        const no = emptyNotification(type)
        no.id = i.toString()
        no.user = demoUsers(1)[0]
        if (Math.random() > 0.5) {
          // like ink
          no.subjectType = SubjectType.Ink
          no.subject = demoInks(1)[0]
        } else {
          no.subjectType = SubjectType.Comment
          no.subject = demoComments(1)[0]
        }
        likes.push(no)
      }
      return likes
    case NotificationType.Reply:
      const replies = []
      n = n == undefined ? 20 : n
      for (let i = 0; i < n; i++) {
        const no = emptyNotification(type)
        no.id = i.toString()
        no.user = demoUsers(1)[0]
        if (Math.random() > 0.5) {
          // reply ink
          no.subjectType = SubjectType.Ink
          no.subject = demoInks(1)[0]
          no.content = {
            commentId: 1,
            source: demoCommentPayload(),
            target: null,
          }
        } else {
          // reply comment
          no.subjectType = SubjectType.Comment
          no.subject = demoComments(1)[0]
          no.content = {
            commentId: 1,
            source: demoCommentPayload(),
            target: demoComments(1)[0].payload,
          }
        }
        replies.push(no)
      }
      return replies
    case NotificationType.System:
      const systems = []
      n = n == undefined ? 20 : n
      for (let i = 0; i < n; i++) {
        const no = emptyNotification(type)
        no.id = i.toString()
        no.user = null
        no.subjectType = SubjectType.System
        no.subject = null
        no.content = 'hello, this is a system notification'
        systems.push(no)
      }
      return systems
  }
}

export const demoMergedLike = (n?: number) => {
  const res = []
  n = n == undefined ? 20 : n
  for (let i = 0; i < n; i++) {
    const no = emptyMergedLike(Math.random() > 0.5)
    no.users = demoUsers((Math.ceil(Math.random() * 10) % 2) + 1)
    if (no.subjectType == SubjectType.Ink) {
      no.subject = demoInks(1)[0]
    } else {
      no.subject = demoComments(1)[0]
    }
    no.total = no.users.length > 1 ? Math.floor(Math.random() * 1000) : 1
    no.updatedAt = randomDate(1)
    res.push(no)
  }
  return res
}

export const demoHotTags = (n?: number) => {
  const tags: TagStats[] = []
  n = n == undefined ? 20 : n
  for (let i = 0; i < n; i++) {
    const tag = {
      name: `tag${i}`,
      reference: Math.floor(Math.random() * 1000),
    }
    tags.push(tag)
  }
  return tags
}
