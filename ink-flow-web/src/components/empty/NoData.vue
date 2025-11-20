<script setup lang="ts">
defineProps({
  image: {
    type: String,
    default: undefined,
  },
  description: {
    type: String,
    default: 'No Data',
  },
})
</script>

<template>
  <el-empty :image-size="200" :image="image" :description="description" class="w-full" />
</template>

<style scoped lang="scss"></style>
