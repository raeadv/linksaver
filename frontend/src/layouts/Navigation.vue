<script setup lang="ts">
import { computed } from 'vue'
import type { MenuOption, NxButton } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/authStore'
import ToggleThemeButton from '@/components/ToggleThemeButton.vue'
import { NMenu } from 'naive-ui'
import { appUtils } from '@/utils/app.utils'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const { dialog, message } = appUtils()
const store = useAuthStore()

const menuOptions: MenuOption[] = [
    { label: 'Home', key: '/' },
    { label: 'Login', key: '/login' },
    { label: 'Register', key: '/register' },
    { label: 'Links', key: '/links' },
    { label: 'Logout', key: 'logout', invoke: true },
];

const activeKey = computed(() => route.path)

const onMenuSelect = (key: string) => {
    if (key === 'logout') {
        dialog.create({
            closable: true,
            title: 'Logout',
            content: 'Are you sure?',
            positiveText: 'Yes',
            negativeText: 'Cancel',
            onPositiveClick: () => {
                store.clearAuth()
                router.replace('/login')
                message.info("you are logged out")
            }
        })
    }


    if (key !== route.path) {
        void router.push(key)
    }
}

const menus = computed(() =>
    auth.isAuthenticated
        ? menuOptions.filter(m => !['login', 'register'].includes(String(m.label).toLowerCase()))
        : menuOptions.filter(m => ['login', 'register'].includes(String(m.label).toLowerCase()))
)

</script>

<template>
    <div class="desktop-nav flex justify-start">
        <toggle-theme-button />
        <n-menu :value="activeKey" mode="horizontal" :options="menus" responsive @update:value="onMenuSelect" />
    </div>
</template>

<style scoped>
.desktop-nav {
    min-width: 320px;
}
</style>