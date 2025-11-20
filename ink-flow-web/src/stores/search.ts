import { defineStore } from 'pinia'
import { ref } from 'vue'

const searchLogKey = 'search-log'

export const useSearchLog = defineStore('search-log', () => {
  const searchLogs = ref<string[]>([])
  searchLogs.value = localStorage.getItem(searchLogKey)
    ? JSON.parse(localStorage.getItem(searchLogKey) as string)
    : []
  function addSearchLog(log: string) {
    if (searchLogs.value.includes(log)) {
      // 存在， 把最新的排到前面
      searchLogs.value.splice(searchLogs.value.indexOf(log), 1)
      searchLogs.value.unshift(log)
      return
    }
    // 超过10条删除最后一条
    if (searchLogs.value.length >= 10) {
      searchLogs.value.pop()
    }
    searchLogs.value.unshift(log)
    localStorage.setItem(searchLogKey, JSON.stringify(searchLogs.value))
  }

  function clearSearchLogs() {
    searchLogs.value = []
    localStorage.removeItem(searchLogKey)
  }

  function removeSearchLog(index: number) {
    searchLogs.value.splice(index, 1)
    localStorage.setItem(searchLogKey, JSON.stringify(searchLogs.value))
  }
  return {
    searchLogs,
    addSearchLog,
    clearSearchLogs,
    removeSearchLog,
  }
})
