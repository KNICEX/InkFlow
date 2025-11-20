<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { defaultAvatar } from '@/consts/default.ts'
import clsx from 'clsx'
const props = defineProps({
  size: {
    type: Number,
    default: 50,
  },
  src: {
    type: String,
    default: '',
  },
  border: {
    type: Boolean,
    default: false,
  },
  preview: {
    type: Boolean,
    default: false,
  },
  hoverMask: {
    type: Boolean,
    default: false,
  },
})
const avatarClass = computed(() => {
  return clsx({
    'border-solid border-2 border-white': props.border,
    'cursor-pointer rounded-full aspect-square': true,
    'hover-mask': props.hoverMask,
  })
})
const sizeStyle = computed(() => {
  return {
    width: `${props.size}px`,
    height: `${props.size}px`,
  }
})
const src = computed(() => {
  return props.src || defaultAvatar
})
const previewList = computed(() => {
  return props.preview ? [src.value] : []
})
</script>
<template>
  <div :style="sizeStyle">
    <el-image
      :class="avatarClass"
      :style="sizeStyle"
      :src="src"
      fit="cover"
      :preview-src-list="previewList"
      hide-on-click-modal
    />
  </div>
</template>

<style scoped lang="scss">
.hover-mask {
  transition: all 0.2s;
  &:hover {
    filter: brightness(0.8);
  }
}
</style>
