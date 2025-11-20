interface NotificationInfo {
  title?: string
  message?: string
  type?: 'success' | 'warning' | 'info' | 'error'
  duration?: number
}
export const notification = (config: NotificationInfo) => {
  ElNotification({
    title: 'Success',
    type: 'success',
    ...config,
  })
}

export const notifySuccessLogin = () => {
  ElNotification({
    title: 'Success',
    message: '登录成功',
    type: 'success',
  })
}
