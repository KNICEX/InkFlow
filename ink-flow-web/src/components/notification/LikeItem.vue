<script setup lang="ts">
import NotificationItem from '@/components/notification/NotificationItem.vue'
import { type MergedLike, SubjectType } from '@/types/notification.ts'
import type { Comment } from '@/types/comment.ts'
import type { Ink } from '@/types/ink.ts'
import { computed } from 'vue'
import InkReference from '@/components/ink/InkReference.vue'
const props = defineProps({
  mergedLike: {
    type: Object as () => MergedLike<Comment | Ink>,
    required: true,
  },
})

const subjectTypeName = computed(() => {
  return props.mergedLike.subjectType == SubjectType.Comment ? '评论' : '作品'
})
</script>

<template>
  <NotificationItem :users="mergedLike.users" :created-at="mergedLike.updatedAt">
    <template #title>
      <span v-if="mergedLike.users.length == 1"> 赞了你的{{ subjectTypeName }} </span>
      <span v-else>等{{ mergedLike.total }}人赞了你的{{ subjectTypeName }}</span>
    </template>
    <div
      v-if="mergedLike.subjectType == SubjectType.Comment"
      class="text-gray-600 mt-2 line-clamp-3"
    >
      <div v-if="mergedLike.subject != null">
        {{ (mergedLike.subject as Comment).payload.content }}
      </div>
      <div v-else>评论已删除</div>
    </div>
    <div v-else class="mt-2">
      <InkReference :ink="mergedLike.subject as Ink"></InkReference>
    </div>
  </NotificationItem>
</template>

<style scoped lang="scss"></style>
