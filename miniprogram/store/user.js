import { defineStore } from 'pinia'
import { getUserInfo } from '../api/auth'
import {
  bootstrapAuth,
  clearTokens,
  getAccessToken,
  getUser,
  isLoggedIn,
  setTokens as setTokensStorage,
  setUser as setUserStorage,
  toLoginPage,
} from '../utils/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    accessToken: getAccessToken(),
    user: getUser(),
    loading: false,
  }),
  getters: {
    loggedIn: (state) => !!state.accessToken,
  },
  actions: {
    syncToken() {
      this.accessToken = getAccessToken()
    },
    setTokens(payload) {
      setTokensStorage(payload || {})
      this.syncToken()
    },
    setUser(user) {
      this.user = user || null
      setUserStorage(user || null)
    },
    async bootstrap() {
      this.loading = true
      try {
        await bootstrapAuth()
        this.syncToken()
        if (!isLoggedIn()) {
          this.setUser(null)
          return
        }
        try {
          const user = await getUserInfo()
          this.setUser(user)
        } catch {
          // ignore
        }
      } finally {
        this.loading = false
      }
    },
    async refreshUser() {
      if (!isLoggedIn()) return null
      const user = await getUserInfo()
      this.setUser(user)
      return user
    },
    logout() {
      clearTokens()
      this.setUser(null)
      toLoginPage()
    },
  },
})

