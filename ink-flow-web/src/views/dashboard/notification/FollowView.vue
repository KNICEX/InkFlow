<script setup lang="ts">
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import FollowList from '@/components/notification/FollowList.vue'
import type { Notification } from '@/types/notification.ts'
import { onMounted, ref } from 'vue'
import { followNotification } from '@/service/notification.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { useUnreadNotificationStore } from '@/stores/notification.ts'

const follows = ref<Notification<never, never>[]>([])
const noStore = useUnreadNotificationStore()
const limit = 15
const { loadMore, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  const res = await followNotification({
    maxId,
    limit,
  })
  if (res.length == 0) {
    return '0'
  }
  follows.value = [...follows.value, ...res]
  return follows.value[follows.value.length - 1].id
})
onMounted(() => {
  loadMore().then(() => {
    noStore.unreadMap['follow'] = 0
  })
})
</script>

<template>
  <DashboardContent title="Follow" class="overflow-hidden">
    <FollowList
      class="overflow-y-auto"
      :follows="follows"
      :load-more="loadMore"
      :loading="loading"
    ></FollowList>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
