<script setup lang="ts">
import LikeItem from '@/components/notification/LikeItem.vue'
import type { MergedLike } from '@/types/notification.ts'
import type { Comment } from '@/types/comment.ts'
import type { Ink } from '@/types/ink.ts'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  mergedLikes: {
    type: Array as () => MergedLike<Comment | Ink>[],
    default: () => [],
  },
  loadMore: {
    type: Function,
    default: () => {},
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const active = useActive()
</script>

<template>
  <div v-infinite-scroll="loadMore" :infinite-scroll-disabled="!active">
    <LikeItem
      class="mb-2"
      v-for="(like, idx) in mergedLikes"
      :mergedLike="like"
      :key="idx"
    ></LikeItem>
    <InkLoading v-show="loading"></InkLoading>
    <NoData v-show="!loading && mergedLikes.length == 0"></NoData>
  </div>
</template>

<style scoped lang="scss"></style>
