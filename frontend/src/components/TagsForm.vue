<script setup lang="js">
import api from '@/net/api';
import { appUtils } from '@/utils/app.utils';
import { Add } from '@vicons/ionicons5';
import { NButton, NCard, NFormItem, NInput, NModal } from 'naive-ui';
import { ref } from 'vue';
import { ta } from 'zod/v4/locales';

const showModal = ref(false)
const tag = ref('')
const { message, loading } = appUtils()

const addTag = async () => {
    if (!tag.value || String(tag.value).length < 1) {
        return
    }

    if (String(tag.value) > 150) {
        message.error("max tag name is 100 chars")
        return
    }

    const { status } = await api.post('/api/tags', { tag: tag.value })

    if (status === 201) {
        message.success(`new Tag ${tag} added`)
        tag.value = ''
        showModal.value = false
    } else {
        message.error(`failed to add tag`)
        return
    }

}

</script>

<template>
    <NButton type="primary" @click="showModal = true">
        <span>New Tags</span>
        <template #icon>
            <Add />
        </template>
    </NButton>


    <NModal v-model:show="showModal" style="max-width: 300px;">
        <NCard title="Add New Tag">
            <div class="flex flex-start space-x-2 justify-around items-center">
                <NFormItem label="Tag Name">
                    <NInput v-model:value="tag" />
                </NFormItem>
                <NButton type="primary" @click="addTag">
                    <span>Add Tag</span>
                    <template #icon>
                        <Add />
                    </template>
                </NButton>
            </div>

        </NCard>
    </NModal>
</template>