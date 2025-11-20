<script setup lang="ts">
import UserList from '@/components/list/user/UserList.vue'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { searchUser } from '@/service/search.ts'
import { ref, watch } from 'vue'
import type { User } from '@/types/user.ts'
import { useProvideFollowHandler } from '@/hook/follow.ts'

const props = defineProps({
  keyword: {
    type: String,
    default: '',
  },
})

const users = ref<User[]>([])
const limit = 15
const { loadMore, reset, loading } = wrapOffsetPagedFunc(async (offset: number) => {
  console.log('load more user', offset)
  if (props.keyword == '') {
    return
  }
  const res = await searchUser({
    keyword: props.keyword,
    offset,
    limit,
  })
  if (res.length == 0) {
    return 0
  }
  users.value = [...users.value, ...res]
  return res.length
}, limit)

useProvideFollowHandler(users)

watch(
  () => props.keyword,

  () => {
    console.log('search-user watch keyword change', props.keyword)
    if (props.keyword == '') {
      return
    }
    reset()
    users.value = []
    loadMore()
  },
)
</script>

<template>
  <div>
    <UserList :load-more="loadMore" :loading="loading" :users="users"></UserList>
  </div>
</template>

<style scoped lang="scss"></style>
