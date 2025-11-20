<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { profile } from '@/service/user.ts'
import type { User } from '@/types/user.ts'
import UserAvatar from '@/components/UserAvatar.vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user.ts'
import InkDialog from '@/components/InkDialog.vue'
import UserInfoEditor from '@/views/user/UserInfoEditor.vue'
import { defaultBanner } from '@/consts/default.ts'
import FollowButton from '@/components/FollowButton.vue'
import { useProvideFollowHandler } from '@/hook/follow.ts'

const userStore = useUserStore()
const route = useRoute()
const router = useRouter()
const userInfo = ref<User | null>(null)

const activeNav = ref('')

const showEdit = ref(false)
const handleEdit = () => {
  showEdit.value = true
}

watch(
  activeNav,
  () => {
    if (!activeNav.value) {
      return
    }
    if (route.name == activeNav.value) {
      return
    }
    router.push({
      name: activeNav.value,
    })
  },
  { immediate: true },
)

const followUserUsers = ref<User[]>([])
useProvideFollowHandler(followUserUsers)

watch(
  () => route.params.account,
  () => {
    if (route.params.account == undefined) {
      return
    }
    activeNav.value = route.name as string
    const account = route.params.account as string
    console.log('watch route params account', route.params.account)
    profile({
      account,
    }).then((val) => {
      userInfo.value = val
      followUserUsers.value = [userInfo.value]
    })
  },
  { immediate: true },
)

const banner = computed(() => {
  if (userInfo.value?.banner == '') {
    return defaultBanner
  }
  return userInfo.value?.banner
})

const isSelf = computed(() => {
  return userStore.getActiveUser()?.user.id == userInfo.value?.id
})
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <div class="w-full relative flex justify-center">
      <el-image
        :src="banner"
        fit="cover"
        class="absolute w-full h-100"
        hide-on-click-modal
        :preview-src-list="[banner]"
      ></el-image>
      <div class="max-screen-w absolute bottom-10 flex justify-between line-padding">
        <div class="flex justify-between items-center">
          <UserAvatar
            :size="80"
            class="mr-6"
            :src="userInfo?.avatar"
            :border="true"
            :preview="true"
          ></UserAvatar>
          <div>
            <div class="text-3xl text-white">{{ userInfo?.username }}</div>
            <div class="text-bg text-white mt-2">@{{ userInfo?.account }}</div>
          </div>
        </div>
        <div>
          <div>
            <el-button v-if="isSelf" size="large" round @click="handleEdit"
              >编辑个人资料
            </el-button>
            <FollowButton
              v-else
              :uid="userInfo?.id ?? '0'"
              size="large"
              :followed="userInfo?.followed"
            ></FollowButton>
          </div>
          <div></div>
        </div>
      </div>
    </div>
    <div class="line-padding max-screen-w">
      <div class="relative text-lg sticky-top white-bg z-1">
        <el-tabs v-model="activeNav">
          <el-tab-pane label="最新" name="latest"></el-tab-pane>
          <el-tab-pane label="点赞" name="likes" v-if="isSelf"></el-tab-pane>
          <!--        <el-tab-pane label="收藏" name="collection"></el-tab-pane>-->
          <el-tab-pane label="浏览" name="views" v-if="isSelf"></el-tab-pane>
          <el-tab-pane label="收藏" name="favorites" v-if="isSelf"></el-tab-pane>
        </el-tabs>
        <div class="absolute right-10 top-3 flex label-text-color">
          <div
            @click="activeNav = 'following'"
            :class="[
              'cursor-pointer hover:text-[var(--primary-color)]',
              activeNav == 'following' ? 'active-color' : '',
            ]"
          >
            正在关注{{ userInfo?.following }}
          </div>
          <div
            @click="activeNav = 'followers'"
            :class="[
              'ml-8 cursor-pointer hover:text-[var(--primary-color)]',
              activeNav == 'followers' ? 'active-color' : '',
            ]"
          >
            关注者{{ userInfo?.followers }}
          </div>
        </div>
      </div>
      <div class="flex justify-center mt-4">
        <router-view v-slot="{ Component }">
          <keep-alive :max="5">
            <Component :is="Component" :uid="userInfo?.id ?? 0"> </Component>
          </keep-alive>
        </router-view>
        <!--        <router-view :uid="userInfo?.id"></router-view>-->
      </div>
    </div>
    <InkDialog v-model="showEdit" title="编辑个人资料">
      <UserInfoEditor :user="userInfo!"></UserInfoEditor>
    </InkDialog>
  </div>
</template>

<style scoped lang="scss">
.active-color {
  color: var(--primary-color);
}
</style>
