<script setup lang="ts">
import CommentSearchList from '@/components/list/comment/CommentSearchList.vue'
import { wrapOffsetPagedFunc } from '@/utils/pagedLoadWrap.ts'
import { searchComment } from '@/service/search.ts'
import { ref, watch } from 'vue'
import type { Comment } from '@/types/comment.ts'

const props = defineProps({
  keyword: {
    type: String,
    default: '',
  },
})

const comments = ref<Comment[]>([])
const limit = 15
const { loadMore, reset, loading } = wrapOffsetPagedFunc(async (offset: number) => {
  console.log('load more comment', offset)
  if (props.keyword == '') {
    return
  }
  const res = await searchComment({
    keyword: props.keyword,
    offset,
    limit,
  })
  if (res.length == 0) {
    return 0
  }
  comments.value = [...comments.value, ...res]
  return res.length
}, limit)

watch(
  () => props.keyword,
  () => {
    if (props.keyword == '') {
      return
    }
    reset()
    comments.value = []
    loadMore()
  },
)
</script>

<template>
  <div>
    <CommentSearchList
      :comments="comments"
      :loading="loading"
      :load-more="loadMore"
    ></CommentSearchList>
  </div>
</template>

<style scoped lang="scss"></style>
