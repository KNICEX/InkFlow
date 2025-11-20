<script setup lang="ts">
import { emptyInteractive, type Interactive } from '@/types/ink.ts'

const props = defineProps({
  interactive: {
    type: Object as () => Interactive,
    default: () => emptyInteractive(),
  },
})

const emit = defineEmits(['like', 'cancelLike', 'favorite', 'cancelFavorite'])

const handleLike = () => {
  if (props.interactive.liked) {
    emit('cancelLike')
  } else {
    emit('like')
  }
}
const handleFavorite = () => {
  if (props.interactive.favorited) {
    emit('cancelFavorite')
  } else {
    emit('favorite')
  }
}
</script>

<template>
  <div class="flex text-gray-500">
    <el-tooltip content="评论数">
      <div class="mr-8 flex items-center cursor-pointer hover:text-[var(--primary-color)]">
        <span class="material-symbols-outlined mr-1 small-icon"> chat </span>
        <span class="">{{ interactive.commentCnt ?? 0 }}</span>
      </div>
    </el-tooltip>
    <el-tooltip content="点赞数">
      <div
        :class="[
          'mr-8 flex items-center cursor-pointer hover:text-red-400',
          interactive.liked ? 'text-red-400' : '',
        ]"
        @click.stop="handleLike"
        title="点赞"
      >
        <span class="material-symbols-outlined mr-1 small-icon"> favorite </span>
        <span class="">{{ interactive.likeCnt ?? 0 }}</span>
      </div>
    </el-tooltip>
    <el-tooltip content="收藏数">
      <div
        :class="[
          'mr-8 flex items-center cursor-pointer hover:text-[var(--primary-color)]',
          interactive.favorited ? 'text-[var(--primary-color)]' : '',
        ]"
        @click.stop="handleFavorite"
        title="收藏"
      >
        <span class="material-symbols-outlined mr-1 small-icon"> bookmark </span>
        <span class="">{{ interactive.favoriteCnt ?? 0 }}</span>
      </div>
    </el-tooltip>
    <el-tooltip content="浏览量">
      <div
        class="flex items-center cursor-pointer hover:text-[var(--primary-color)]"
        title="浏览量"
      >
        <span class="material-symbols-outlined mr-1 small-icon"> bar_chart_4_bars </span>
        <span class="">{{ interactive.viewCnt ?? 0 }}</span>
      </div>
    </el-tooltip>
  </div>
</template>

<style scoped lang="scss">
.small-icon {
  font-size: 1.3rem;
}
</style>
