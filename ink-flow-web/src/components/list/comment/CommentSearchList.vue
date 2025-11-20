<script setup lang="ts">
import { demoComments } from '@/mock/demo_data.ts'
import CommentItem from '@/components/list/comment/CommentItem.vue'
import type { Comment } from '@/types/comment.ts'
import NoData from '@/components/empty/NoData.vue'
import InkLoading from '@/components/loading/InkLoading.vue'
import BackTop from '@/components/BackTop.vue'
import { useActive } from '@/hook/active.ts'
defineProps({
  comments: {
    type: Array as () => Comment[],
    default: () => demoComments(10, false),
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
  <div>
    <div v-infinite-scroll="loadMore" :infinite-scroll-disabled="!active">
      <CommentItem
        v-for="comment in comments"
        :comment="comment"
        :key="comment.id"
        :reply-able="false"
      ></CommentItem>
    </div>
    <InkLoading v-show="loading"></InkLoading>
    <BackTop></BackTop>
    <NoData v-show="!loading && comments.length == 0"></NoData>
  </div>
</template>

<style scoped lang="scss"></style>
