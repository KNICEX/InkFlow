<script setup lang="ts">
import { ref, useTemplateRef } from 'vue'
import { ElPopover } from 'element-plus'

const props = defineProps({
  round: {
    type: Boolean,
    default: true,
  },
  place: {
    type: String as () =>
      | 'bottom-start'
      | 'bottom'
      | 'bottom-end'
      | 'left'
      | 'right'
      | 'top'
      | 'top-start'
      | 'top-end',
    default: 'bottom-start',
  },
  trigger: {
    type: String as () => 'hover' | 'click' | 'focus' | '',
    default: 'hover',
  },
  padding: {
    type: String,
    default: '',
  },
  showAfter: {
    type: Number,
    default: 10,
  },
  width: {
    type: String,
    default: 'auto',
  },
})

const popStyle = ref({
  borderRadius: '1rem',
  padding: props.padding,
})

const popoverRef = useTemplateRef<InstanceType<typeof ElPopover>>('popoverRef')

defineExpose({
  hide: () => {
    popoverRef.value?.hide()
  },
})
</script>

<template>
  <el-popover
    ref="popoverRef"
    :width="width"
    :show-after="showAfter"
    :trigger="trigger == '' ? undefined : trigger"
    :popper-style="popStyle"
    popper-class=" overflow-hidden
 px-0 shadow-xl round-popover rounded-xl "
    :placement="place"
    :show-arrow="false"
    effect="light"
  >
    <template #default>
      <slot name="content"></slot>
    </template>
    <template #reference>
      <slot name="reference"></slot>
    </template>
  </el-popover>
</template>

<style scoped lang="scss"></style>
