<script setup lang="ts">
import { provide, type Ref, ref } from 'vue'
const props = defineProps({
  defaultActive: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['on-change'])
const activeItem = ref(props.defaultActive)

const updateActiveItem = (key: string) => {
  activeItem.value = key
  emit('on-change', key)
}

export interface ActiveItemInjection {
  activeItem: Ref<string>
  updateActiveItem: (key: string) => void
}

provide('active-item', { activeItem, updateActiveItem })
const setActive = (key: string) => {
  updateActiveItem(key)
}
defineExpose({
  setActive,
})
</script>

<template>
  <div
    class="flex flex-col w-full items-start px-4 h-full"
    v-for="(slot, name) in $slots"
    :key="name"
  >
    <slot :name="name"></slot>
  </div>
</template>

<style scoped lang="scss"></style>
