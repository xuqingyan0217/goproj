import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userId = ref<string | null>(localStorage.getItem('userId'))
  const cartNum = ref(parseInt(localStorage.getItem('cartNum') || '0'))
  const error = ref('')
  const warning = ref('')

  const setUserId = (id: string | null) => {
    userId.value = id
    if (id) {
      localStorage.setItem('userId', id)
    } else {
      localStorage.removeItem('userId')
    }
  }

  const setCartNum = (num: number) => {
    cartNum.value = num
    localStorage.setItem('cartNum', num.toString())
  }

  const setError = (msg: string) => {
    error.value = msg
  }

  const setWarning = (msg: string) => {
    warning.value = msg
  }

  const logout = async () => {
    try {
      await fetch('/auth/logout', { method: 'POST' })
      setUserId(null)
      setCartNum(0)
      localStorage.clear()
      window.location.href = '/sign-in'
    } catch (err) {
      error.value = 'Logout failed'
    }
  }

  return {
    userId,
    cartNum,
    error,
    warning,
    setUserId,
    setCartNum,
    setError,
    setWarning,
    logout
  }
})