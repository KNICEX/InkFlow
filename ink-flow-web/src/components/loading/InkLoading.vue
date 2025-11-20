<script setup lang="ts">
defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  text: {
    type: String,
    default: 'Loading...',
  },
})
</script>

<template>
  <div v-loading="true" class="w-full h-10"></div>
</template>

<style scoped lang="scss"></style>
