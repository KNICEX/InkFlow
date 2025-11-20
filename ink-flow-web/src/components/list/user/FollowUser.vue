<script setup lang="ts">
import ImageInk from '@/components/list/ink/ImageInk.vue'
import UserCard from '@/components/UserCard.vue'
import type { UserFollow } from '@/types/user.ts'

defineProps({
  user: {
    type: Object as () => UserFollow,
    default: () => ({}),
  },
})
</script>

<template>
  <div class="flex h-74 pt-2 border-gray-200 border-b-1">
    <div class="w-90 pt-2">
      <UserCard :show-follow="false" :user="user"></UserCard>
    </div>
    <div class="overflow-x-hidden">
      <div ref="scrollDiv" class="flex overflow-x-auto h-full items-start pt-4 no-scroll pr-10">
        <div class="w-50 mr-5 flex-shrink-0" v-for="ink in user.inks" :key="ink.id">
          <ImageInk :ink="ink" :show-author="false"></ImageInk>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.mask {
  position: relative;
}

.no-scroll {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.button {
  position: absolute;
  height: 100%;
  width: 6rem;
  //background: rgba(128, 128, 128, 0.31);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;

  span {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #141414;
    height: 3rem;
    width: 3rem;
    border-radius: 50%;
    opacity: 30%;
  }
}

.left-button {
  left: 0;
  top: 0;
  background: linear-gradient(to right, white, rgba(255, 255, 255, 0));
  span {
    padding-left: 0.4rem;
  }
}

.right-button {
  right: 0;
  top: 0;
  background: linear-gradient(to left, white, rgba(255, 255, 255, 0));
}
</style>
