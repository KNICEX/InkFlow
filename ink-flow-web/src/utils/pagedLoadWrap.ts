import { ref } from 'vue'
import type { Ref } from 'vue'
type offsetLoadFunc = (offset: number) => Promise<number | undefined>
type maxIdLoadFunc = (maxId: string) => Promise<string | undefined>
type maxIdTimestampLoadFunc = (
  maxId: string,
  timestamp: number,
) => Promise<{ maxId: string; timestamp: Date } | undefined>

interface WrappedFunc {
  loadMore: () => void
  reset: () => void
  loading: Ref<boolean>
}

export const wrapOffsetPagedFunc = (loadFunc: offsetLoadFunc, limit: number) => {
  let offset = 0
  const loading = ref(false)
  let noMore = false

  // 限流用
  let lastLoad = new Date(0)
  let lastOffset = 0
  let lastReqHash: string | undefined = undefined
  return {
    loadMore: async (reqHash?: string) => {
      // 限流
      const now = new Date()
      if (
        now.getTime() - lastLoad.getTime() < 100 &&
        lastOffset == offset &&
        lastReqHash === reqHash
      ) {
        return
      } else {
        lastLoad = now
        lastOffset = offset
        lastReqHash = reqHash
      }

      if (loading.value) return
      if (noMore) return
      loading.value = true
      const n = await loadFunc(offset)
      if (n === undefined) {
        // 不需要实际请求，可能是限流或者其他原因
        loading.value = false
        return true
      }

      offset += n
      loading.value = false
      if (n < limit) {
        noMore = true
      }
      return !noMore // false if no more data
    },
    reset: () => {
      offset = 0
      noMore = false
      loading.value = false
    },
    loading,
  }
}

export const wrapMaxIdPagedFunc = (loadFunc: maxIdLoadFunc) => {
  let maxId = '0'
  const loading = ref(false)
  let noMore = false

  // 限流用
  let lastLoad = new Date(0)
  let lastMaxId = '0'

  return {
    loadMore: async () => {
      // 限流
      const now = new Date()
      if (now.getTime() - lastLoad.getTime() < 100 && lastMaxId == maxId) {
        return
      } else {
        lastLoad = now
        lastMaxId = maxId
      }

      if (loading.value) return
      if (noMore) {
        loading.value = false
        return
      }
      loading.value = true
      const resMaxId = await loadFunc(maxId)
      if (resMaxId === undefined) {
        // 不需要实际请求，可能是限流或者其他原因
        loading.value = false
        return true
      }

      if (resMaxId == '0') {
        loading.value = false
        noMore = true
        return false // no more data
      }
      maxId = resMaxId
      loading.value = false
      return !noMore // false if no more data
    },
    reset: () => {
      maxId = '0'
      noMore = false
      loading.value = false
    },
    loading,
  }
}

export const wrapMaxIdTimestampPagedFunc = (loadFunc: maxIdTimestampLoadFunc) => {
  let maxId = '0'
  let timestamp = new Date()
  const loading = ref(false)
  let noMore = false
  return {
    loadMore: async () => {
      if (loading.value) return
      if (noMore) return
      loading.value = true
      const resMaxIdTimestamp = await loadFunc(maxId, timestamp.getTime())
      if (resMaxIdTimestamp === undefined) {
        // 不需要实际请求，可能是限流或者其他原因
        loading.value = false
        return true
      }

      if (resMaxIdTimestamp.maxId == '0') {
        loading.value = false
        noMore = true
        return false // no more data
      }
      maxId = resMaxIdTimestamp.maxId
      timestamp = new Date(resMaxIdTimestamp.timestamp)
      loading.value = false
      return !noMore // false if no more data
    },
    reset: () => {
      maxId = '0'
      timestamp = new Date()
      noMore = false
      loading.value = false
    },
    loading,
  }
}
