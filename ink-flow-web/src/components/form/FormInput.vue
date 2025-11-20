<script setup lang="ts">
const model = defineModel()
defineProps({
  type: {
    type: String,
    default: 'text',
  },
  validateEvent: {
    type: Boolean,
    default: true,
  },
  maxlength: {
    type: Number,
    default: 100,
  },
  clearable: {
    type: Boolean,
    default: true,
  },
})
</script>

<template>
  <el-input
    v-model="model"
    :maxlength="maxlength"
    :validate-event="validateEvent"
    :show-password="type == 'password'"
    :type="type"
    class="rounded-sm-input h-12"
    :clearable="clearable"
    rows="3"
    resize="none"
  >
    <template #append v-if="$slots.append">
      <slot name="append"></slot>
    </template>
  </el-input>
</template>

<style scoped lang="scss">
.el-form-item__content {
  height: 1rem;
}
</style>
