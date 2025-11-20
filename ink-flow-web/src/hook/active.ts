import { onActivated, onDeactivated, ref } from 'vue'

export const useActive = () => {
  const active = ref(true)
  onActivated(() => {
    active.value = true
  })
  onDeactivated(() => {
    active.value = false
  })
  return active
}
