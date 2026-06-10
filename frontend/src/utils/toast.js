const createToast = (message, type = 'info', duration = 3000) => {
  const bgColors = {
    success: 'bg-green-500',
    error: 'bg-red-500',
    warning: 'bg-yellow-500',
    info: 'bg-blue-500'
  }

  const toastEl = document.createElement('div')
  toastEl.className = `fixed top-4 left-1/2 transform -translate-x-1/2 z-50 px-6 py-3 rounded-lg shadow-lg text-white font-medium ${bgColors[type]}`
  toastEl.textContent = message

  document.body.appendChild(toastEl)

  setTimeout(() => {
    toastEl.style.opacity = '0'
    toastEl.style.transition = 'opacity 0.3s ease'
    setTimeout(() => {
      document.body.removeChild(toastEl)
    }, 300)
  }, duration)
}

const Toast = {
  success: (msg) => createToast(msg, 'success'),
  error: (msg) => createToast(msg, 'error'),
  warning: (msg) => createToast(msg, 'warning'),
  info: (msg) => createToast(msg, 'info')
}

export default Toast
