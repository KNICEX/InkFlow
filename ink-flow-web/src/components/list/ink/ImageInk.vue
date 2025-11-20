<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import UserCard from '@/components/UserCard.vue'
import InkPopover from '@/components/popover/InkPopover.vue'
import { type Ink, inkStatusProp } from '@/types/ink.ts'
import { useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import { useProvideFollowHandler } from '@/hook/follow.ts'
import type { User } from '@/types/user.ts'
const router = useRouter()
const props = defineProps({
  ink: {
    type: Object as () => Ink,
    required: true,
  },
  showAuthor: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['on-cover-click'])
const handleCoverClick = () => {
  console.log('cover click')
  emit('on-cover-click', props.ink.id)
}
const handleToUserHome = () => {
  router.push(`/user/${props.ink.author.account}`)
}

const toDetailHref = computed(() => {
  return `/ink/${props.ink.id}?status=${inkStatusProp(props.ink.status)}`
})
useProvideFollowHandler(ref<User[]>([props.ink.author]))
</script>

<template>
  <div class="rounded-1 cursor-pointer">
    <div class="flex flex-col justify-between">
      <router-link
        :to="toDetailHref"
        class="overflow-hidden rounded-xl border border-gray-200 dark:border-gray-600 aspect-square"
      >
        <el-image
          v-if="ink.cover != ''"
          @click="handleCoverClick"
          class="w-full aspect-square hover:scale-105 transition duration-300"
          loading="lazy"
          fit="cover"
          :src="ink.cover"
        >
        </el-image>
        <div
          v-else
          @click="handleCoverClick"
          class="flex justify-center hover:scale-115 transition duration-300 text-xl items-center w-full aspect-square"
        >
          <div class="px-8 line-clamp-2">{{ ink.title }}</div>
        </div>
      </router-link>
    </div>
    <div class="flex flex-col mt-2" @click="handleToUserHome">
      <router-link :to="toDetailHref" class="line-clamp-2">{{ ink.title }}</router-link>
      <div class="flex" v-if="showAuthor">
        <InkPopover :show-after="300">
          <template #reference>
            <div class="flex items-center mt-2">
              <UserAvatar class="mr-3" :size="25" :src="ink.author.avatar"></UserAvatar>
              <span>{{ ink.author.username }}</span>
            </div>
          </template>
          <template #content>
            <UserCard :user="ink.author"></UserCard>
          </template>
        </InkPopover>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
