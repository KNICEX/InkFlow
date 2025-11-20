<script setup lang="ts">
import { ref, watch } from 'vue'
import CommentList from '@/components/list/comment/CommentList.vue'
import { type Comment, type CommentPayload, emptyCommentStats } from '@/types/comment.ts'
import { provide } from 'vue'
import {
  likeComment,
  loadComment,
  reply,
  deleteComment,
  cancelLikeComment,
} from '@/service/comment.ts'
import { useUserStore } from '@/stores/user.ts'
import { wrapMaxIdPagedFunc } from '@/utils/pagedLoadWrap.ts'

const props = defineProps({
  biz: {
    type: String,
    default: 'ink',
  },
  bizId: {
    type: String,
    default: '0',
  },
})

const userStore = useUserStore()

const limit = 15
const comments = ref<Comment[]>([])
const { loadMore, reset, loading } = wrapMaxIdPagedFunc(async (maxId: string) => {
  if (props.bizId == '0') {
    return
  }
  const res = await loadComment({
    maxId,
    biz: props.biz,
    bizId: props.bizId,
    limit,
  })

  if (res.length == 0) {
    return '0'
  }
  comments.value = [...comments.value, ...res]
  return comments.value[comments.value.length - 1].id
})
watch(
  () => props.bizId,
  () => {
    reset()
    comments.value = []
    loadMore()
  },
  { immediate: true },
)

const handleLike = async (c: Comment) => {
  await likeComment(c.id)
  c.stats.liked = true
  c.stats.likeCnt += 1
}
const handleCancelLike = async (c: Comment) => {
  await cancelLikeComment(c.id)
  c.stats.liked = false
  c.stats.likeCnt -= 1
}

provide('handleCommentInteractive', {
  handleLike,
  handleCancelLike,
})

const handleReply = async (parent: Comment | null, payload: CommentPayload) => {
  const id = await reply({
    biz: props.biz,
    bizId: props.bizId,
    parentId: parent?.id ?? '0',
    rootId: parent?.root?.id ?? parent?.id ?? '0',
    payload: payload,
  })
  const c: Comment = {
    id,
    biz: props.biz,
    bizId: props.bizId,
    payload,
    commentator: userStore.getActiveUser()?.user!,
    isAuthor: false, // TODO这里暂时好像还不知道
    parent: parent,
    root: parent?.root ?? null,
    children: [],
    stats: emptyCommentStats(),
    createdAt: new Date(),
  }
  console.log('parent', parent)
  if (parent) {
    if (parent.children) {
      parent.children.push(c)
    } else {
      parent.children = [c]
    }
  } else {
    comments.value.unshift(c)
  }
}

const handleDelete = async (c: Comment) => {
  await deleteComment(c.id)
  if (!c.root) {
    // 一级评论， 直接删除
    for (let i = 0; i < comments.value.length; i++) {
      if (comments.value[i].id == c.id) {
        comments.value.splice(i, 1)
        break
      }
    }
    return
  }
  for (let i = 0; i < comments.value.length; i++) {
    // 二级评论， 找到对应根评论
    const comment = comments.value.find((v) => v.id == c?.root?.id)
    if (!comment || !comment.children) {
      return
    }
    for (let j = 0; j < comment.children.length; j++) {
      if (comment.children[j]!.id == c.id) {
        comment.children.splice(j, 1)
        break
      }
    }
  }
}
</script>

<template>
  <div>
    <CommentList
      class="mt-6"
      :loading="loading"
      :comments="comments"
      @submit-comment="handleReply"
      @delete-comment="handleDelete"
    ></CommentList>
  </div>
</template>

<style scoped lang="scss"></style>
