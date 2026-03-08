<script setup lang="js">
import LinkForm from '@/components/LinkForm.vue';
import TagsForm from '@/components/TagsForm.vue';
import UserInfo from '@/components/UserInfo.vue';
import { useAuthStore } from '@/store/authStore';
import { computed, ref } from 'vue';
import Login from './auth/Login.vue';
import LinksData from '@/components/LinksData.vue';

const store = useAuthStore()
const linksDataRef = ref()

const handleRefreshLinkData = () => {
    linksDataRef.value.fetchLinks()
}
</script>

<template>
    <div v-if="store.isAuthenticated" class="w-full">
        <UserInfo />
        <div class="mt-2 my-2 w-full flex flex-row space-x-2">
            <TagsForm />
            <LinkForm @link-created="handleRefreshLinkData" />
        </div>
        <div class="mt-2">
            <LinksData ref="linksDataRef" />
        </div>
    </div>
    <div v-else class="w-full flex items-center justify-center">
        <Login />
    </div>
</template>
