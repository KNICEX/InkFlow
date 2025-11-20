<script setup lang="ts">
import { computed, ref } from 'vue'
import ImageList from '@/components/upload/ImageList.vue'

defineProps({
  placeholder: {
    type: String,
    default: '写下你的评论...',
  },
  showCancel: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['submit', 'cancel'])

const commentContent = ref('')
const isSubmitting = ref(false)
// 存储上传的图片URL
const images = ref<{ url: string }[]>([])
const imageUrls = computed(() => {
  return images.value.map((image) => image.url)
})

const submitComment = () => {
  if (!commentContent.value.trim() && images.value.length === 0) return

  isSubmitting.value = true

  // 提交评论，增加images参数
  emit('submit', {
    content: commentContent.value,
    images: imageUrls.value,
  })

  // 重置
  commentContent.value = ''
  images.value = []
  isSubmitting.value = false
}

const handleCancel = () => {
  commentContent.value = ''
  images.value = []
  emit('cancel')
}
</script>

<template>
  <div class="rounded-xl p-4 border border-gray-200 dark:border-gray-700">
    <el-input
      v-model="commentContent"
      :placeholder="placeholder"
      type="textarea"
      :rows="3"
      resize="none"
      class="mb-3"
    />
    <div class="flex justify-between items-start">
      <ImageList list-type="picture" v-model="images">
        <el-tooltip content="添加图片" placement="top">
          <el-button type="primary" text circle>
            <span class="material-symbols-outlined">image</span>
          </el-button>
        </el-tooltip>
      </ImageList>

      <div class="flex">
        <el-button v-if="showCancel" @click="handleCancel" class="mr-2">取消</el-button>
        <el-button
          type="primary"
          @click="submitComment"
          :loading="isSubmitting"
          :disabled="!commentContent.trim() && images.length === 0"
        >
          发布评论
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
:deep(.el-textarea__inner) {
  background-color: #f9f9f9;
  border-color: #e0e0e0;

  .dark & {
    background-color: #2d3748;
    border-color: #4a5568;
    color: #e2e8f0;
  }
}

.hidden {
  display: none;
}
</style>
