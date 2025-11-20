<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import UserCard from '@/components/UserCard.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import type { Ink } from '@/types/ink.ts'
import { useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import InkInteractive from '@/components/ink/InkInteractive.vue'
import { formatDate } from '@/utils/date.ts'
import { useInjectInkInteractiveHandler } from '@/hook/interactive.ts'
import { useProvideFollowHandler } from '@/hook/follow.ts'
import type { User } from '@/types/user.ts'

const router = useRouter()
const props = defineProps({
  ink: {
    type: Object as () => Ink,
    required: true,
  },
  showAuthor: {
    type: Boolean,
    default: true,
  },
  previewLength: {
    type: Number,
    default: 150,
  },
})

const contentPreview = computed(() => {
  // 从文章内容中提取文本预览
  if (!props.ink.contentHtml) return ''
  // const text = props.ink.contentHtml.replace(/<[^>]*>/g, '') // 简单去除HTML标签
  // return text.length > props.previewLength ? text.substring(0, props.previewLength) + '...' : text
  // 去除html，保留换行
  const text = props.ink.contentHtml.replace(/<[^>]*>/g, '\n')
  return text.length > props.previewLength ? text.substring(0, props.previewLength) + '...' : text
  // return text
})

const emit = defineEmits(['on-cover-click'])
const handleCoverClick = () => {
  emit('on-cover-click', props.ink.id)
}

const handleItemClick = () => {
  router.push('/ink/' + props.ink.id)
}

const handleToUserHome = () => {
  router.push(`/user/${props.ink.author.account}`)
}

const formattedDate = computed(() => {
  // 格式化日期，这里简单返回，实际应用中可能需要使用日期格式化函数
  return formatDate(props.ink.createdAt)
})

const { handleLike, handleFavorite, handleCancelLike, handleCancelFavorite } =
  useInjectInkInteractiveHandler()

useProvideFollowHandler(ref<User[]>([props.ink.author]))
</script>

<template>
  <div class="flex rounded-border p-4 mb-4 cursor-pointer hover-bg" @click="handleItemClick">
    <div class="mr-3">
      <InkPopover place="bottom" :show-after="300">
        <template #reference>
          <div class="flex items-center">
            <UserAvatar :size="48" :src="ink.author.avatar"></UserAvatar>
          </div>
        </template>
        <template #content>
          <UserCard :user="ink.author"></UserCard>
        </template>
      </InkPopover>
    </div>
    <!-- 作者信息 -->
    <div>
      <div class="flex items-center" v-if="showAuthor">
        <div class="flex" @click.stop="handleToUserHome">
          <span class="font-semibold">{{ ink.author.username }}</span>
          <el-link class="text-gray-500 ml-2">@{{ ink.author.account }}</el-link>
          <span class="text-gray-500 ml-2 flex items-end">{{ formattedDate }}</span>
        </div>
      </div>

      <!-- 标题 -->
      <!--      <h3 class="text-lg">{{ ink.title }}</h3>-->

      <!-- 内容预览 -->
      <p class="mb-3">{{ contentPreview }}</p>

      <!-- 封面图片 -->
      <div v-if="ink.cover" class="mb-3 flex">
        <el-image
          preview-teleported
          @click.stop="handleCoverClick"
          fit="contain"
          class="ink-image transition duration-300"
          loading="lazy"
          :z-index="10000"
          hide-on-click-modal
          :preview-src-list="[ink.cover]"
          :src="ink.cover"
        >
        </el-image>
      </div>

      <!-- 标签 -->
      <div class="flex flex-wrap gap-2 mt-2" v-if="ink.tags && ink.tags.length">
        <el-tag v-for="tag in ink.tags" :key="tag" size="small" type="info">#{{ tag }}</el-tag>
      </div>

      <!-- 互动数据 -->
      <InkInteractive
        class="mt-2"
        :interactive="ink.interactive"
        @like="handleLike(ink.id)"
        @cancel-like="handleCancelLike(ink.id)"
        @favorite="handleFavorite(ink.id)"
        @cancel-favorite="handleCancelFavorite(ink.id)"
      ></InkInteractive>
    </div>
  </div>
</template>

<style scoped lang="scss">
.ink-image {
  max-height: 30rem;
  max-width: 100%;
  width: auto;
  height: auto;
  &:deep(.el-image__inner) {
    width: auto;
    height: 100%;
    border-radius: 0.8rem;
  }
}
</style>
