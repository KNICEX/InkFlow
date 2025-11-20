interface confirmOptions {
  title?: string
  message?: string
  type?: 'success' | 'warning' | 'info' | 'error'
  confirmed: () => void
  canceled?: () => void
}
export const confirm = (opts: confirmOptions) => {
  ElMessageBox.confirm(opts.message ?? '你确定要继续操作吗😨?', opts.title ?? 'Warning', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: opts.type ?? 'warning',
    confirmButtonClass: 'rounded-xl',
  })
    .then(async () => {
      opts.confirmed()
    })
    .catch(() => {
      opts.canceled?.()
    })
}
