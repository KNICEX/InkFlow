<script setup lang="ts">
import { emptyInk, type Ink } from '@/types/ink.ts'
import { formatDate } from '@/utils/date.ts'
import InkCover from '@/components/image/InkCover.vue'

const props = defineProps({
  ink: {
    type: Object as () => Ink,
    required: true,
  },
})

let checkInk = props.ink
if (checkInk == null) {
  checkInk = emptyInk()
  checkInk.title = '帖子不存在'
  checkInk.cover = 'https://images.app.goo.gl/mfYr8ih6wDddXAZo7'
}
</script>

<template>
  <router-link :to="checkInk.id != 0 ? `/ink/${checkInk.id}` : ''">
    <div class="flex">
      <InkCover
        class="w-40 h-30 rounded-xl"
        :src="checkInk.cover"
        :backoff="checkInk.title"
      ></InkCover>
      <div class="flex flex-col justify-between ml-6">
        <div class="nav-text">
          {{ checkInk.title }}
        </div>
        <div>
          <el-tag type="info" class="mr-2" v-for="(tag, idx) in checkInk.tags" :key="idx">{{
            tag
          }}</el-tag>
        </div>
        <div class="text-sm text-gray-500 mt-1">{{ formatDate(checkInk.createdAt) }}</div>
      </div>
    </div>
  </router-link>
</template>

<style scoped lang="scss"></style>
