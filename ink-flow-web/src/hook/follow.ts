import { provide, inject } from 'vue'
import type { Ref } from 'vue'
import type { User } from '@/types/user.ts'
import { cancelFollow, follow } from '@/service/user.ts'
import { notification } from '@/utils/notification.ts'

interface config {
  before?: () => boolean
  onError?: (error: Error) => void
}

export const useProvideFollowHandler = (users: Ref<User[]>) => {
  const handleFollow = async (uid: string) => {
    console.log('handleFollow', uid)
    await follow(uid)
    const user = users.value.find((u) => u.id === uid)
    if (user) {
      user.followed = true
      user.followers++
    }
    notification({
      message: '关注成功',
    })
  }
  const handleCancelFollow = async (uid: string) => {
    await cancelFollow(uid)
    const user = users.value.find((u) => u.id === uid)
    if (user) {
      user.followed = false
      user.followers--
    }
    notification({
      message: '取消关注成功',
    })
  }
  provide('handleFollow', handleFollow)
  provide('handleCancelFollow', handleCancelFollow)
}

export const useInjectFollowHandler = () => {
  const doNothing = () => {}
  const handleFollow = inject<(uid: string) => void>('handleFollow') ?? doNothing
  const handleCancelFollow = inject<(uid: string) => void>('handleCancelFollow') ?? doNothing
  return {
    handleFollow,
    handleCancelFollow,
  }
}
