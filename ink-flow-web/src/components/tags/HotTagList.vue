<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { topTag } from '@/service/stats'
import type { TagStats } from '@/types/stats'

const loading = ref(false)
const hotTags = ref<TagStats[]>([])
onMounted(async () => {
  hotTags.value = await topTag(0, 5)
})
</script>

<template>
  <div class="rounded-border">
    <div class="text-lg font-bold px-4 pt-4 pb-2">热门标签</div>
    <el-skeleton :loading="loading" animated :count="5">
      <template #template>
        <div class="flex items-center justify-between mb-3">
          <el-skeleton-item variant="text" />
        </div>
      </template>
      <template #default>
        <router-link
          :to="`/search/ink?key=${encodeURIComponent(`#${tag.name}`)}`"
          v-for="tag in hotTags"
          :key="tag.name"
          class="hover-bg py-3 cursor-pointer px-4 flex flex-col"
        >
          <div class="semibold-text">{{ tag.name }}</div>
          <div class="text-xs text-gray-500 mt-1">{{ tag.reference }} 条帖子</div>
        </router-link>
      </template>
    </el-skeleton>
  </div>
</template>

<style scoped lang="scss"></style>
