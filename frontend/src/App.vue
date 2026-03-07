<script setup>
import { NConfigProvider, NModalProvider } from 'naive-ui'
import { RouterView } from 'vue-router';
import { uiStore } from '@/store/uiStore';
import { computed, watchEffect } from 'vue';
import { darkTheme } from 'naive-ui'

import Navigation from './layouts/Navigation.vue';

const ui = uiStore()
const providerTheme = computed(() => (ui.theme === 'dark' ? darkTheme : null))

watchEffect(() => {
  document.documentElement.dataset.theme = ui.theme
})

</script>

<template>
  <n-config-provider :theme="providerTheme">
    <NModalProvider>
      <nav>
        <Navigation />
      </nav>
      <main class="w-full flex flex-col px-2 justify-center">
        <RouterView />
      </main>
    </NModalProvider>
  </n-config-provider>
</template>

<style scoped></style>
