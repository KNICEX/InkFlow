<script setup lang="ts">
import { computed, inject } from 'vue'
import UserAvatar from '@/components/UserAvatar.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import type { Comment } from '@/types/comment.ts'
import MoreOperation from '@/components/button/MoreOperation.vue'
import UserCard from '@/components/UserCard.vue'
import { formatDate } from '@/utils/date.ts'
import { useUserStore } from '@/stores/user.ts'
import { confirm } from '@/utils/message.ts'

const props = defineProps({
  comment: {
    type: Object as () => Comment,
    required: true,
  },
  isChild: {
    type: Boolean,
    default: false,
  },
  replyAble: {
    type: Boolean,
    default: true,
  },
  showStats: {
    type: Boolean,
    default: true,
  },
})

const userStore = useUserStore()

const emit = defineEmits(['reply', 'delete', 'like', 'cancel-like'])

const handleReply = () => {
  emit('reply', props.comment)
}

const { handleLike: doLike, handleCancelLike: doCancelLike } = inject<{
  handleLike: (comment: Comment) => Promise<void>
  handleCancelLike: (comment: Comment) => Promise<void>
}>('handleCommentInteractive') ?? {
  handleLike: async () => {},
  handleCancelLike: async () => {},
}
const handleLike = () => {
  if (props.comment.stats.liked) {
    console.log('cancel like')
    doCancelLike(props.comment)
    emit('cancel-like', props.comment)
    return
  }
  doLike(props.comment)
  emit('like', props.comment)
}

const handleDelete = () => {
  confirm({
    message: '你确定要删除此条回复吗😨?',
    confirmed: () => {
      emit('delete', props.comment)
    },
  })
}

const likeIconClass = computed(() => {
  return props.comment.stats.liked ? 'text-red-500' : 'text-gray-500'
})

const operations = (c: Comment) => {
  const commonOps = [
    {
      name: '复制',
      action: () => {
        navigator.clipboard.writeText(c.payload.content)
      },
    },
    {
      name: '举报',
      action: () => {
        console.log('举报')
      },
    },
  ]
  if (c.commentator.id == userStore.getActiveUser()?.user.id) {
    commonOps.unshift({
      name: '删除',
      action: () => {
        confirm({
          message: '你确定要删除此条回复吗😨?',
          confirmed: () => {
            emit('delete', c)
          },
        })
      },
    })
  }
  return commonOps
}
</script>

<template>
  <div class="flex p-4 rounded-xl transition-colors duration-200 hover-bg">
    <InkPopover :show-after="300" place="bottom">
      <template #reference>
        <UserAvatar :hover-mask="true" :src="comment.commentator.avatar" :size="46" />
      </template>
      <template #content>
        <UserCard :user="comment.commentator"></UserCard>
      </template>
    </InkPopover>
    <div class="w-full ml-4">
      <div class="flex items-start">
        <div class="flex-grow">
          <div class="flex justify-between">
            <div class="text-lg">
              <router-link :to="`/user/${comment.commentator.account}`" class="text-lg">
                {{ comment.commentator.username }}
                <el-link class="text-gray-500 text-lg">@{{ comment.commentator.account }} </el-link>
              </router-link>

              <span class="ml-3 text-sm text-gray-500">{{ formatDate(comment.createdAt) }}</span>
            </div>

            <MoreOperation :horizon="true" :operations="operations(comment)"></MoreOperation>
          </div>
        </div>
      </div>

      <div>
        <div
          v-if="isChild && comment.parent && comment.parent.id != comment.root?.id"
          class="text-gray-500 mb-2"
        >
          回复
          <router-link
            class="text-[var(--primary-color)] hover:underline"
            :to="`/user/${comment.parent.commentator.account}`"
          >
            @{{ comment.parent.commentator.username }}
          </router-link>
        </div>

        <div class="leading-relaxed mb-3 break-words">{{ comment.payload.content }}</div>

        <div
          v-if="comment.payload.images && comment.payload.images.length > 0"
          class="flex flex-wrap gap-2 mt-3"
        >
          <div v-for="(image, index) in comment.payload.images" :key="index" class="relative">
            <el-image
              :src="image"
              fit="cover"
              :preview-src-list="comment.payload.images"
              hide-on-click-modal
              :initial-index="index"
              preview-teleported
              class="w-28 h-28 rounded-lg object-cover cursor-pointer border border-gray-200 dark:border-gray-600 transition-transform duration-200 hover:scale-105"
            />
          </div>
        </div>
      </div>

      <div class="flex mt-4" v-if="showStats">
        <div
          class="flex items-center mr-6 cursor-pointer text-gray-500 hover:text-blue-500"
          @click="handleReply()"
        >
          <span class="material-symbols-outlined text-xl">chat_bubble_outline</span>
          <span v-if="comment.children && comment.stats.replyCnt > 0" class="text-sm ml-1">
            {{ comment.stats.replyCnt }}
          </span>
        </div>

        <div
          class="flex items-center mr-6 cursor-pointer hover:text-red-500"
          :class="likeIconClass"
          @click="handleLike"
        >
          <span class="material-symbols-outlined text-xl">favorite</span>
          <span v-if="comment.stats.likeCnt > 0" class="text-sm ml-1">
            {{ comment.stats.likeCnt }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
