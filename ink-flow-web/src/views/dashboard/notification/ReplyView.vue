<script setup lang="ts">
import { onMounted, ref } from 'vue'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import ReplyList from '@/components/notification/ReplyList.vue'
import type { Notification, ReplyContent } from '@/types/notification.ts'
import { replyNotification } from '@/service/notification.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { useUnreadNotificationStore } from '@/stores/notification.ts'

const replies = ref<Notification<never, ReplyContent>[]>([])
const noStore = useUnreadNotificationStore()
const limit = 15
const { loadMore, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  const res = await replyNotification({
    maxId,
    limit,
  })
  if (res.length == 0) {
    return '0'
  }
  replies.value = [...replies.value, ...res]
  if (res.length < limit) {
    return '0'
  }
  return replies.value[replies.value.length - 1].id
})
onMounted(() => {
  loadMore().then(() => {
    noStore.unreadMap['reply'] = 0
  })
})
</script>

<template>
  <DashboardContent title="Reply" class="overflow-hidden">
    <ReplyList class="overflow-y-auto" :replies="replies" :load-more="loadMore" :loading="loading">
    </ReplyList>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
