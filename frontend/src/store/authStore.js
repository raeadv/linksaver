import { defineStore } from "pinia";
import { jwtDecode } from "jwt-decode";

export const useAuthStore = defineStore("auth", {
    state: () => ({
        access_token: localStorage.getItem("access_token") || "",
        expires_at: localStorage.getItem("expires_at") || "",
        user: {
            username: null,
            email: null
        },
    }),
    getters: {
        isAuthenticated: (s) => !!s.access_token,
        userdata: (s) => {
            if (!s.access_token) return null

            const parsed = jwtDecode(s.access_token)
            return { username: parsed.username, email: parsed.email }
        },
    },
    actions: {
        setLogin(token) {
            const parsed = jwtDecode(token)

            console.log({ token, parsed })
            if (parsed) {
                this.access_token = token;
                this.expires_at = parsed.expires_at;
                this.user = parsed.user;
                localStorage.setItem("access_token", token);
                localStorage.setItem("expires_at", parsed.expires_at);
            }
        },

        clearAuth() {
            this.access_token = "";
            this.expires_at = ""
            this.user = null;
            localStorage.removeItem("access_token");
            localStorage.removeItem("expires_at");
        },
    },
});
