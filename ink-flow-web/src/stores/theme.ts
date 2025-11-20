import { ref, nextTick } from 'vue'
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref('light')
  function loadFromLocalStorage() {
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      changeTheme(savedTheme)
    }
  }

  function saveToLocalStorage() {
    localStorage.setItem('theme', theme.value)
  }
  function changeTheme(val: string) {
    theme.value = val
    saveToLocalStorage()
    if (theme.value == 'dark') {
      nextTick(() => {
        document.documentElement.classList.add('dark')
      })
    } else {
      nextTick(() => {
        document.documentElement.classList.remove('dark')
      })
    }
  }

  function toDark() {
    changeTheme('dark')
  }

  function toLight() {
    changeTheme('light')
  }

  return {
    theme,
    loadFromLocalStorage,
    saveToLocalStorage,
    toDark,
    toLight,
  }
})
