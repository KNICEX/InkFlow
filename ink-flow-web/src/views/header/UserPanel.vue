<script setup lang="ts">
import { emptyUser } from '@/types/user.ts'
import InkPopover from '@/components/popover/InkPopover.vue'
import InkAvatar from '@/components/UserAvatar.vue'
import { ArrowDownBold } from '@element-plus/icons-vue'
import UserPop from '@/views/header/UserPop.vue'
import { useTemplateRef } from 'vue'
import { useUserStore } from '@/stores/user.ts'
const userStore = useUserStore()
const userPanelPopRef = useTemplateRef<InstanceType<typeof InkPopover>>('userPanelPopRef')
const hide = () => {
  userPanelPopRef.value?.hide()
}
</script>

<template>
  <InkPopover ref="userPanelPopRef" place="bottom" trigger="click" padding="0">
    <template #reference>
      <div
        class="flex items-center justify-between w-60 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl px-4 py-2 cursor-pointer"
      >
        <div class="flex items-center">
          <InkAvatar :src="userStore.getActiveUser()?.user.avatar"></InkAvatar>
          <div class="ml-2">
            <div class="semibold-text">
              {{ userStore.getActiveUser()?.user.username }}
            </div>
            <div class="text-sm text-gray-400">@{{ userStore.getActiveUser()?.user.account }}</div>
          </div>
        </div>
        <el-icon class="ml-3">
          <ArrowDownBold class="text-gray-400" />
        </el-icon>
      </div>
    </template>
    <template #content>
      <UserPop :close="hide" :user="userStore.getActiveUser()?.user ?? emptyUser()"></UserPop>
    </template>
  </InkPopover>
</template>

<style scoped lang="scss"></style>
