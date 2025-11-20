<script setup lang="ts">
defineProps({
  src: {
    type: String,
    default: '',
  },
  alt: {
    type: String,
    default: '',
  },
  backoff: {
    type: String,
    default: 'cover',
  },
})
</script>

<template>
  <div class="overflow-hidden">
    <el-image
      v-if="src != ''"
      :alt="alt"
      :fit="'cover'"
      class="w-full h-full"
      :src="src"
    ></el-image>
    <div v-else class="w-full h-full flex transition-all justify-center items-center text-xl px-2">
      <div class="line-clamp-1">{{ backoff }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
