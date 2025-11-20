<script setup lang="ts">
import { ref, defineModel } from 'vue'
import type { UploadProps } from 'element-plus'
import { uploadImage } from '@/service/file.ts'

defineProps({
  listType: {
    type: String as () => 'picture' | 'picture-card' | 'text',
    default: undefined,
  },
})

const images = defineModel<
  {
    url: string
    name?: string
  }[]
>({
  default: () => [],
})

const dialogVisible = ref(false)
const dialogImageUrl = ref('')
const preview: UploadProps['onPreview'] = (file) => {
  dialogImageUrl.value = file.url!
  dialogVisible.value = true
}
const handleUpload: UploadProps['beforeUpload'] = async (image) => {
  const url = await uploadImage(image)
  if (url) {
    images.value.push({
      url: url,
      name: image.name,
    })
  }
  return false
}
</script>

<template>
  <el-upload
    v-model:file-list="images"
    action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15"
    :list-type="listType"
    multiple
    :on-preview="preview"
    :before-upload="handleUpload"
  >
    <slot>
      <span class="material-symbols-outlined"> add </span>
    </slot>
  </el-upload>

  <el-dialog v-model="dialogVisible">
    <img w-full :src="dialogImageUrl" alt="Preview Image" />
  </el-dialog>
</template>

<style scoped lang="scss"></style>
