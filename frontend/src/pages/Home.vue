<script setup lang="js">
import LinkForm from '@/components/LinkForm.vue';
import TagsForm from '@/components/TagsForm.vue';
import UserInfo from '@/components/UserInfo.vue';
import { useAuthStore } from '@/store/authStore';
import { computed } from 'vue';

const store = useAuthStore()
const user = computed(() => store.userdata)

</script>

<template>
    <UserInfo />
    <div class="mt-2 my-2 w-full flex flex-row space-x-2">
        <TagsForm />
        <LinkForm />
    </div>
</template>
