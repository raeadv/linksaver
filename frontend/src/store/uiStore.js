import { defineStore } from "pinia";


const parseTheme = (value) => (value === "dark" ? "dark" : "light");

export const uiStore = defineStore("ui", {
    state: () => ({
        theme: parseTheme(localStorage.getItem("theme")),
    }),
    getters: {
        currentTheme: (ui) => ui.theme
    },
    actions: {
        boot() {
            this.theme = parseTheme(localStorage.getItem("theme"));
        },
        toggleTheme() {
            this.theme = this.theme === "light" ? "dark" : "light";
            localStorage.setItem('theme', this.theme)
        }
    }
});