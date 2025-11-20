<script setup lang="ts">
import { computed, ref } from 'vue'
import { useInjectFollowHandler } from '@/hook/follow.ts'
import { useShowLoginStore, useUserStore } from '@/stores/user.ts'

const userStore = useUserStore()
const showLoginStore = useShowLoginStore()
const hasLogin = computed(() => {
  return !!userStore.getActiveUser()?.user.id
})

defineProps({
  uid: {
    type: String,
    required: true,
  },
  size: {
    type: String,
    default: 'default',
  },
  followed: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['follow', 'cancelFollow'])

const { handleFollow: injectFollow, handleCancelFollow: injectCancelFollow } =
  useInjectFollowHandler()

const followTip = ref('正在关注')
const followedColor = ref('')
const handleMouseover = () => {
  followTip.value = '取消关注'
  followedColor.value = 'red'
}
const handleMouseLeave = () => {
  followTip.value = '正在关注'
  followedColor.value = ''
}

const handleFollow = (uid: string) => {
  if (!hasLogin.value) {
    showLoginStore.showLogin = true
    return
  }
  emit('follow', uid)
  injectFollow(uid)
}

const handleCancelFollow = (uid: string) => {
  emit('cancelFollow', uid)
  injectCancelFollow(uid)
}
</script>

<template>
  <el-button
    v-if="followed"
    :color="followedColor"
    @click="handleCancelFollow(uid)"
    @mouseover="handleMouseover"
    @mouseleave="handleMouseLeave"
    :size="size"
    plain
    round
  >
    {{ followTip }}
  </el-button>
  <el-button v-else color="black" :size="size" round @click="handleFollow(uid)"> 关注 </el-button>
</template>

<style scoped lang="scss"></style>
