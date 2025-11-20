<script setup lang="ts">
import { ref, watch } from 'vue'
import type { UserFollow } from '@/types/user.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { following } from '@/service/relation.ts'
import { useProvideFollowHandler } from '@/hook/follow.ts'
import FollowUserList from '@/components/list/user/FollowUserList.vue'

const props = defineProps({
  uid: {
    type: String,
    default: '0',
  },
})

const limit = 15
const users = ref<UserFollow[]>([])

const { loadMore, reset, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  if (props.uid == '0') {
    return
  }

  const res = await following(props.uid, { maxId, limit })
  if (res.length == 0) {
    return '0'
  }
  users.value = [...users.value, ...res]
  if (res.length < limit) {
    return '0'
  }
  return users.value[users.value.length - 1].id
})

watch(
  () => props.uid,
  () => {
    console.log('following uid change', props.uid)
    reset()
    users.value = []
    loadMore()
  },
  { immediate: true },
)

useProvideFollowHandler(users)
</script>

<template>
  <div class="w-full">
    <FollowUserList :users="users" :load-more="loadMore" :loading="loading"></FollowUserList>
  </div>
</template>

<style scoped lang="scss"></style>
