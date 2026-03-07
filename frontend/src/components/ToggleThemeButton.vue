<script setup lang="ts">
import { uiStore } from '@/store/uiStore';
import { computed } from 'vue';
import { Sunny, Moon } from '@vicons/ionicons5'
import { NButton, NIcon } from 'naive-ui'

const ui = uiStore()
const isDarkTheme = computed(() => ui.theme === 'dark')
</script>

<template>

    <span v-if="isDarkTheme">
        <n-button quaternary circle @click="ui.toggleTheme">
            <template #icon>
                <n-icon>
                    <Sunny />
                </n-icon>
            </template>
        </n-button>
    </span>
    <span v-else>
        <n-button quaternary circle @click="ui.toggleTheme">
            <template #icon>
                <n-icon>
                    <Moon />
                </n-icon>
            </template>
        </n-button>
    </span>

</template>