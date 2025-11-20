<script setup lang="ts">
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import { ref, useTemplateRef } from 'vue'

const props = defineProps({
  src: {
    type: String, // url or base64
    default: '',
  },
  imageName: {
    type: String,
    default: '',
  },
  aspectRatio: {
    type: Number,
    default: undefined,
  },
})

const emit = defineEmits(['confirm', 'cancel'])

const cropperRef = useTemplateRef<InstanceType<typeof Cropper>>('cropperRef')

export interface CropperResult {
  src: string
  name: string
  canvas: HTMLCanvasElement
}
const handleConfirm = () => {
  if (cropperRef.value) {
    const { image, canvas } = cropperRef.value.getResult()
    emit('confirm', {
      src: image.src,
      name: props.imageName,
      canvas: canvas,
    })
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>
<template>
  <div class="image-cropper">
    <Cropper
      ref="cropperRef"
      :src="src"
      :stencil-props="{
        aspectRatio: aspectRatio,
      }"
    />
    <div class="cropper-actions mt-4 flex justify-end gap-2">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确认</el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.image-cropper {
  width: 100%;
}
</style>
