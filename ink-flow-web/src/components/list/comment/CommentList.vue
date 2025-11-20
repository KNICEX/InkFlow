<script setup lang="ts">
import { type PropType, ref } from 'vue'
import CommentItem from './CommentItem.vue'
import type { Comment, CommentPayload } from '@/types/comment.ts'
import NoData from '@/components/empty/NoData.vue'
import CommentInput from '@/components/list/comment/CommentInput.vue'
import { useActive } from '@/hook/active.ts'

defineProps({
  comments: {
    type: Array as () => Comment[],
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  loadMore: {
    type: Function,
    default: () => {},
  },
  loadMoreChild: {
    type: Function as PropType<(c: Comment) => void>,
    default: () => {},
  },
})

const replyingComment = ref<Comment | null>(null)

const emit = defineEmits(['submit-comment', 'delete-comment'])

const showReply = (comment: Comment) => {
  replyingComment.value = comment
}

const handleSubmitComment = (parentComment: Comment | null = null, payload: CommentPayload) => {
  emit('submit-comment', parentComment, payload)
  if (parentComment) {
    replyingComment.value = null
  }
}

const cancelReply = () => {
  replyingComment.value = null
}

const handleDelete = (c: Comment) => {
  emit('delete-comment', c)
}

const active = useActive()
</script>

<template>
  <div class="w-full" v-infinite-scroll="loadMore" :infinite-scroll-disabled="!active">
    <div class="mb-4">
      <CommentInput
        placeholder="发表你的评论..."
        @submit="(payload) => handleSubmitComment(null, payload)"
      />
    </div>

    <div v-if="comments.length === 0" class="py-8">
      <NoData description="没有更多评论"></NoData>
    </div>
    <div v-else>
      <div v-for="comment in comments" :key="comment.id" class="mb-2 last:mb-0">
        <CommentItem :comment="comment" @reply="showReply" @delete="handleDelete" />

        <div v-if="replyingComment && replyingComment.id === comment.id" class="ml-12 mt-2 mb-4">
          <CommentInput
            :placeholder="`回复 @${comment.commentator.username}...`"
            @submit="(payload) => handleSubmitComment(comment, payload)"
            @cancel="cancelReply"
            :show-cancel="true"
          />
        </div>

        <div
          v-if="comment.children && comment.children.length > 0"
          class="ml-12 mt-2 pl-4 border-l-2 border-gray-200 dark:border-gray-700"
        >
          <div
            v-for="childComment in comment.children"
            :key="childComment.id"
            class="mb-2 last:mb-0"
          >
            <CommentItem
              :comment="childComment"
              :is-child="true"
              :show-child="false"
              @reply="showReply"
              @delete="handleDelete"
            />

            <div v-if="replyingComment && replyingComment.id === childComment.id" class="mt-2 mb-4">
              <CommentInput
                :placeholder="`回复 @${childComment.commentator.username}...`"
                @submit="(payload) => handleSubmitComment(childComment, payload)"
                @cancel="cancelReply"
                :show-cancel="true"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="py-4">
      <el-skeleton :rows="3" animated />
    </div>
  </div>
</template>
<style lang="scss"></style>
