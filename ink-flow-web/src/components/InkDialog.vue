<script setup lang="ts">
import { ArrowLeftBold, CloseBold } from '@element-plus/icons-vue'
const model = defineModel()
defineProps({
  title: {
    type: String,
    default: 'TITLE',
  },
  showClose: {
    type: Boolean,
    default: true,
  },
  showBack: {
    type: Boolean,
    default: false,
  },
  width: {
    type: String,
    default: '30rem',
  },
  height: {
    type: String,
    default: 'auto',
  },
})

const emit = defineEmits(['onClose', 'onBack'])
const onClose = () => {
  emit('onClose')
}
const onBack = () => {
  emit('onBack')
}
</script>
<template>
  <el-dialog
    ref="dialogRef"
    append-to="#app"
    v-model="model"
    @closed="onClose"
    :style="{ height: height }"
    :show-close="false"
    :width="width"
    class="rounded-dialog"
  >
    <template #header>
      <slot name="header">
        <div class="flex justify-between items-center c">
          <slot name="header-before">
            <div
              v-if="showBack"
              @click="onBack"
              class="cursor-pointer p-2 text-gray-400 flex justify-center rounded-sm items-center hover:bg-gray-100"
            >
              <el-icon :size="20"><ArrowLeftBold /></el-icon>
            </div>
            <div v-else></div>
          </slot>
          <h2 class="font-semibold">{{ title }}</h2>
          <div>
            <div
              v-if="showClose"
              @click="model = false"
              class="cursor-pointer p-2 text-gray-400 flex justify-center rounded-sm items-center hover:bg-gray-100"
            >
              <el-icon :size="20">
                <CloseBold />
              </el-icon>
            </div>
          </div>
        </div>
      </slot>
    </template>
    <slot name="default"> This is a content </slot>
    <slot name="footer"> </slot>
  </el-dialog>
</template>

<style scoped lang="scss"></style>
