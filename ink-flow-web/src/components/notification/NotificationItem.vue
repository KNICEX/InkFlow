<script setup lang="ts">
import InkPopover from '@/components/popover/InkPopover.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import UserCard from '@/components/UserCard.vue'
import type { User } from '@/types/user.ts'
import MoreOperation from '@/components/button/MoreOperation.vue'
import { useProvideFollowHandler } from '@/hook/follow.ts'
import { ref } from 'vue'
import { formatDate } from '../../utils/date.ts'

const props = defineProps({
  users: {
    type: Array as () => User[],
    required: true,
  },
  system: {
    type: Boolean,
    default: false,
  },
  createdAt: {
    type: Date,
  },
})

const handleDelNo = () => {}

const ops = [
  {
    name: '删除',
    action: handleDelNo,
  },
]

useProvideFollowHandler(ref(props.users))
</script>

<template>
  <div class="flex max-w-240 hover:bg-gray-100 p-4 rounded-xl transition-colors cursor-pointer">
    <template v-if="!system">
      <div v-if="users.length == 1" class="w-20">
        <InkPopover :show-after="400">
          <template #reference>
            <UserAvatar :hover-mask="true" :src="users[0].avatar"></UserAvatar>
          </template>
          <template #content>
            <UserCard :user="users[0]"></UserCard>
          </template>
        </InkPopover>
      </div>
      <div v-else class="relative w-20">
        <InkPopover :show-after="400">
          <template #reference>
            <UserAvatar :hover-mask="true" :size="40" :src="users[0].avatar"></UserAvatar>
          </template>
          <template #content>
            <UserCard :user="users[0]"></UserCard>
          </template>
        </InkPopover>
        <InkPopover :show-after="400">
          <template #reference>
            <UserAvatar
              :hover-mask="true"
              :size="40"
              :src="users[1].avatar"
              class="absolute left-4 top-4"
              border
            ></UserAvatar>
          </template>
          <template #content>
            <UserCard :user="users[1]"></UserCard>
          </template>
        </InkPopover>
      </div>
    </template>
    <template v-else>
      <div class="w-20">
        <span class="text-4xl">🤖</span>
      </div>
    </template>
    <div class="w-full">
      <div class="flex justify-between">
        <div>
          <div v-if="!system" class="flex">
            <div v-if="users.length == 1">
              <span>
                <router-link :to="`/user/${users[0].account}`">
                  <span>{{ users[0].username }}</span>
                  <el-link class="nav-text ml-2">@{{ users[0].account }}</el-link>
                </router-link>
              </span>
              <span class="ml-2">
                <slot name="title"></slot>
              </span>
            </div>
            <div v-else>
              <span>
                <router-link :to="`/user/${users[0].account}`">
                  <span>{{ users[0].username }}</span>
                </router-link>
                <span>、</span>
                <router-link :to="`/user/${users[1].account}`">
                  <span>{{ users[1].username }}</span>
                </router-link>
              </span>
              <span class="ml-4">
                <slot name="title"></slot>
              </span>
            </div>
          </div>
          <div v-else>
            <span>系统通知</span>
          </div>
          <div class="text-sm text-gray-500 mt-1">{{ formatDate(createdAt) }}</div>
        </div>
        <MoreOperation :horizon="true" :operations="ops"></MoreOperation>
      </div>
      <div>
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>
