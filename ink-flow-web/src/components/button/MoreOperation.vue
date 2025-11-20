<script setup lang="ts">
import { computed } from 'vue'
import InkPopover from '@/components/popover/InkPopover.vue'

interface Operation {
  name: string
  icon?: string
  action: () => void
}

const props = defineProps({
  horizon: {
    type: Boolean,
    default: false,
  },
  operations: {
    type: Array as () => Operation[],
    default: () => [],
  },
})

const icon = computed(() => {
  return props.horizon ? 'more_horiz' : 'more_vert'
})
</script>

<template>
  <InkPopover trigger="click" padding="0" place="bottom">
    <template #reference>
      <div
        class="hover:bg-gray-200 cursor-pointer size-10 transition-colors rounded-lg flex justify-center items-center"
      >
        <span class="material-symbols-outlined text-gray-600"> {{ icon }} </span>
      </div>
    </template>
    <template #content>
      <div v-for="(op, idx) in operations" :key="idx" @click="op.action" class="popover-button">
        {{ op.name }}
      </div>
    </template>
  </InkPopover>
</template>

<style scoped lang="scss"></style>
