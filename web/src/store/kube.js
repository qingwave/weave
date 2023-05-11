import { defineStore } from 'pinia'

export const useKubeStore = defineStore('kube', {
  state: () => {
    return { data: {} }
  },
  actions: {
    setNamespace(namespace) {
        this.data.namespace = namespace
    },
    getNamespace() {
        return this.data.namespace
    },
  },
})