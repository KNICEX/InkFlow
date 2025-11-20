import { createRouter, createWebHistory } from 'vue-router'
import { nextTick } from 'vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'layout',
      component: () => import('../views/LayoutView.vue'),
      children: [
        {
          path: '/',
          name: 'home',
          component: () => import('../views/home/HomeView.vue'),
        },
        {
          path: '/search',
          name: 'search',
          component: () => import('../views/search/SearchView.vue'),
          children: [
            {
              path: 'ink',
              name: 'search-ink',
              component: () => import('../views/search/SearchInk.vue'),
            },
            {
              path: 'user',
              name: 'search-user',
              component: () => import('../views/search/SearchUser.vue'),
            },
            {
              path: 'reply',
              name: 'search-reply',
              component: () => import('../views/search/SearchReply.vue'),
            },
          ],
        },
        {
          path: '/user/:account',
          name: 'user',
          component: () => import('../views/user/UserView.vue'),
          children: [
            {
              path: '',
              name: 'latest',
              component: () => import('../views/user/UserInkList.vue'),
            },
            {
              path: 'likes',
              name: 'likes',
              component: () => import('../views/user/InkLiked.vue'),
            },
            {
              path: 'views',
              name: 'views',
              component: () => import('../views/user/InkViewed.vue'),
            },
            {
              path: 'favorites',
              name: 'favorites',
              component: () => import('../views/user/UserFavorites.vue'),
            },
            {
              path: 'following',
              name: 'following',
              component: () => import('../views/user/UserFollowing.vue'),
            },
            {
              path: 'followers',
              name: 'followers',
              component: () => import('../views/user/UserFollowers.vue'),
            },
          ],
        },
        {
          path: '/ink/:id',
          name: 'ink-detail',
          component: () => import('../views/ink/InkDetail.vue'),
        },
      ],
    },

    {
      path: '/dashboard/:account',
      name: 'dashboard',
      component: () => import('../views/dashboard/DashboardLayout.vue'),
      children: [
        {
          path: '',
          name: 'dashboard-home',
          component: () => import('../views/dashboard/DashboardHome.vue'),
        },
        {
          name: 'editor',
          path: 'editor/:id?',
          component: () => import('../views/dashboard/InkEditor.vue'),
        },
        {
          path: 'ink',
          redirect: 'ink/published',
        },
        {
          path: 'ink/:status?',
          name: 'dashboard-ink',
          component: () => import('../views/dashboard/InkView.vue'),
        },
        {
          path: 'system',
          name: 'dashboard-system',
          component: () => import('../views/dashboard/notification/SystemView.vue'),
        },
        {
          path: 'reply',
          name: 'dashboard-reply',
          component: () => import('../views/dashboard/notification/ReplyView.vue'),
        },
        {
          path: 'like',
          name: 'dashboard-like',
          component: () => import('../views/dashboard/notification/LikeView.vue'),
        },
        {
          path: 'follow',
          name: 'dashboard-follow',
          component: () => import('../views/dashboard/notification/FollowView.vue'),
        },
      ],
    },
  ],
})

router.afterEach(() => {
  nextTick(() => {
    window.scrollTo({ top: 0 })
  })
})

export default router
