<script setup lang="ts">
import { type Ink, InkStatus, inkStatusLabel } from '@/types/ink.ts'
import { formatDate } from '@/utils/date.ts'
import InkCover from '@/components/image/InkCover.vue'
import InkPopover from '@/components/popover/InkPopover.vue'

import { computed, inject } from 'vue'
import clsx from 'clsx'
import MoreOperation from '@/components/button/MoreOperation.vue'

const props = defineProps({
  ink: {
    type: Object as () => Ink,
    required: true,
  },
})

const handleInkOp =
  inject<(ink: Ink, op: 'edit' | 'delete' | 'preview' | 'share' | 'publish') => void>(
    'handle-ink-op',
  )
const handleEdit = () => {
  handleInkOp!(props.ink, 'edit')
}
const handlePreview = () => {
  handleInkOp!(props.ink, 'preview')
}
const handleDelete = () => {
  handleInkOp!(props.ink, 'delete')
}
const handleShare = () => {
  console.log('share ink', props.ink.id)
}
const statusClass = computed(() => {
  return clsx({
    'absolute rounded-full w-2.5 h-2.5  right-7.5 top-4': true,
    'bg-green-500': props.ink.status == InkStatus.Published,
    'bg-blue-500': props.ink.status == InkStatus.UnPublished,
    'bg-yellow-500': props.ink.status == InkStatus.Pending,
    'bg-red-500': props.ink.status == InkStatus.Private,
  })
})

const operations = (ink: Ink) => {
  const commonOps = [
    {
      name: '预览',
      action: () => handleInkOp!(ink, 'preview'),
    },
    {
      name: '分享',
      action: () => handleInkOp!(ink, 'share'),
    },
  ]
  if (ink.status != InkStatus.Pending) {
    commonOps.unshift({
      name: '编辑',
      action: () => handleInkOp!(ink, 'edit'),
    })
    commonOps.push({
      name: '删除',
      action: () => handleInkOp!(ink, 'delete'),
    })
  }
  if (ink.status == InkStatus.UnPublished) {
    commonOps.push({
      name: '发布',
      action: () => handleInkOp!(ink, 'publish'),
    })
  }
  return commonOps
}
</script>

<template>
  <div class="flex justify-between max-w-240 cursor-pointer relative">
    <div class="flex">
      <InkCover class="w-40 h-30 rounded-xl" :src="ink.cover" :backoff="ink.title"></InkCover>
      <div class="flex flex-col justify-between ml-8">
        <div class="nav-text">{{ ink.title }}</div>
        <div>
          <span class="mr-1" v-for="tag in ink.tags" :key="tag">
            <el-tag type="info">#{{ tag }}</el-tag>
          </span>
        </div>
        <div class="text-sm text-gray-500">
          <span>{{ inkStatusLabel(ink.status) }}</span>
          <span class="mx-2">•</span>
          <span>{{ formatDate(ink.updatedAt) }}</span>
        </div>
      </div>
    </div>
    <div class="flex items-center">
      <InkPopover place="left" padding="0" trigger="click">
        <template #reference>
          <MoreOperation :horizon="true" :operations="operations(ink)"></MoreOperation>
        </template>
        <template #content>
          <div>
            <div
              v-show="ink.status != InkStatus.Pending"
              class="popover-button"
              @click="handleEdit"
            >
              编辑
            </div>
            <div class="popover-button" @click="handlePreview">预览</div>
            <div
              v-show="ink.status != InkStatus.Pending"
              class="popover-button"
              @click="handleDelete"
            >
              删除
            </div>
            <div class="popover-button" @click="handleShare">分享</div>
          </div>
        </template>
      </InkPopover>
    </div>
    <div :class="statusClass" :title="inkStatusLabel(ink.status)"></div>
  </div>
</template>

<style scoped lang="scss"></style>
