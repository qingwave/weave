import { defineStore } from 'pinia'

export const usePostStore = defineStore('post', {
  state: () => {
    return { data: {} }
  },
  actions: {
    set(val) {
        this.data = val
    },
    get(id) {
      if (this.data.id == id ) {
        return this.data
      }
      return {}
    },
  },
})