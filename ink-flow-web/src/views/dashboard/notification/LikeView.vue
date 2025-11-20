<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { MergedLike } from '@/types/notification.ts'
import DashboardContent from '@/views/dashboard/DashboardContent.vue'
import LikeList from '@/components/notification/LikeList.vue'
import { mergedLikeNotification } from '@/service/notification.ts'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { useUnreadNotificationStore } from '@/stores/notification.ts'

const likes = ref<MergedLike<never>[]>([])
const noStore = useUnreadNotificationStore()

const limit = 15

const { loadMore, loading } = wrapOffsetPagedFunc(async (offset: number) => {
  const res = await mergedLikeNotification({
    offset,
    limit,
  })
  if (res.length == 0) {
    return 0
  }
  likes.value = [...likes.value, ...res]
  return res.length
}, limit)
onMounted(() => {
  loadMore().then(() => {
    noStore.unreadMap['like'] = 0
  })
})
</script>
<template>
  <DashboardContent title="Like" class="overflow-hidden">
    <LikeList
      class="overflow-y-auto"
      :merged-likes="likes"
      :load-more="loadMore"
      :loading="loading"
    ></LikeList>
  </DashboardContent>
</template>

<style scoped lang="scss"></style>
