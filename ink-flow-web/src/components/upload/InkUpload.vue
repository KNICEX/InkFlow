<script setup lang="ts">
import type { UploadRawFile } from 'element-plus'
import { computed } from 'vue'
const props = defineProps({
  multiple: {
    type: Boolean,
    default: false,
  },
  fileType: {
    type: String as () => 'image' | 'any',
    default: 'image',
  },
  drag: {
    type: Boolean,
    default: false,
  },
  max: {
    type: Number,
    default: 1,
  },
  src: {
    type: Array as () => string[],
    default: () => [],
  },
})

const accept = computed(() => {
  return props.fileType === 'image' ? 'image/*' : ''
})
const emit = defineEmits(['upload'])
const handleUpload = (f: UploadRawFile) => {
  emit('upload', f)
  // 阻止默认行为
  return false
}
</script>

<template>
  <el-upload
    multiple
    action=""
    :accept="accept"
    list-type="picture-card"
    :before-upload="handleUpload"
  >
  </el-upload>
</template>

<style scoped lang="scss"></style>
