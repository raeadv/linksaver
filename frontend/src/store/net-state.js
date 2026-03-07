import { defineStore } from "pinia";


export const netState = defineStore("net-state", {
    state: () => ({
        loading: false
    }),
    getters: {
        isLoading: (net) => net.loading
    },
    actions: {
        startLoading() {
            this.loading = true
        },
        stopLoading() {
            this.loading = false
        }
    }
})