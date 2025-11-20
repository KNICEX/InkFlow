<script setup lang="ts">
import { computed, defineProps } from 'vue'
import clsx from 'clsx'

const props = defineProps({
  bg: {
    type: String as () => 'white' | 'gray',
    default: 'white',
  },
  round: {
    type: Boolean,
    default: true,
  },
  border: {
    type: Boolean,
    default: true,
  },
  shadow: {
    type: Boolean,
    default: true,
  },
  padding: {
    type: Boolean,
    default: true,
  },
})

const cardClass = computed(() => {
  return clsx({
    'p-4': props.padding,
    'shadow-md dark:shadow-gray-700': props.shadow,
    'border  p-4  ': props.border,
    'border-gray-100': props.border && props.bg == 'white',
    'border-slate-200 dark:border-slate-600': props.border && props.bg == 'gray',
    'bg-slate-100 dark:bg-gray-600': props.bg === 'gray',
    'rounded-sm': !props.round,
    'rounded-xl': props.round,
  })
})
</script>

<template>
  <div :class="cardClass">
    <slot></slot>
  </div>
</template>

<style scoped lang="scss"></style>
