<script setup lang="js">
import api from '@/net/api';
import { appUtils } from '@/utils/app.utils';
import { createFormRules } from '@/utils/validation.adapter';
import { Add, Search } from '@vicons/ionicons5';
import { NButton, NCard, NForm, NFormItem, NInput, NInputGroup, NModal, NSelect } from 'naive-ui';
import { reactive, ref } from 'vue';
import z from 'zod';

const emit = defineEmits(['link-created'])
const { message, loading } = appUtils()
const showModal = ref(false)
const formRef = ref()
const linkData = reactive({
    link: '',
    name: '',
    link_desc: '',
    link_tags: []
})
const tagsRef = ref([]);

const addLink = async () => {

    const { status, data } = await api.post('/api/links', { ...linkData })
    console.log(data, status)
    if (status === 201) {
        message.success(`new link added`)
        showModal.value = false
        emit('link-created', null)
    } else {
        message.error(`failed to add link`)
        return
    }

}

async function handleSearch(query) {
    if (!query.length || query.length < 2) {
        tagsRef.value = [];
        return;
    }

    const { data, status } = await api.get(`/api/tags/${query}`)

    if (status === 200) {
        tagsRef.value = data.tags || []
    }


}

const handleGetUrlMeta = async () => {
    if (linkData.link.length < 1) {
        return
    }

    message.info('Getting your link metadata...')

    const { data, status } = await api.get('/api/links/get/web-meta', {
        params: {
            url: linkData.link
        }
    })

    if (status !== 200) {
        message.error("Failed to get Website Metadata, you have to set the title and description manually")
    }

    if (status === 200) {
        const metadata = data.metadata

        linkData.name = metadata.title || ''
        linkData.link_desc = metadata.description || ''
    }


}


const schema = z.object({
    link: z.url({ message: "Must be valid url" }),
    name: z.string().min(3, { message: 'minimal url name is 3 chars' }),
    link_desc: z.string().nullable(),
    link_tags: z.array(z.any()).optional()
})


const rules = createFormRules(schema)

const validateAndSaveLink = async (e) => {
    e.preventDefault()
    formRef.value?.validate(async (errors) => {
        if (!errors) {
            await addLink()
        }
        else {
            message.error('save link Failed')
        }
    })
}


</script>

<template>
    <NButton type="primary" @click="showModal = true">
        <span>New Link</span>
        <template #icon>
            <Add />
        </template>
    </NButton>


    <NModal v-model:show="showModal" style="max-width: 300px;">
        <NCard title="Add New Link">
            <NForm ref="formRef" size="large" :label-width="80" :model="linkData" :rules="rules"
                :disabled="loading === true">
                <div class="flex flex-col space-x-2 justify-around items-center">
                    <NFormItem label="Url/Link" path="link">
                        <NInputGroup>
                            <NInput v-model:value="linkData.link" />
                            <NButton @click="handleGetUrlMeta">
                                <template #icon>
                                    <Search />
                                </template>
                            </NButton>
                        </NInputGroup>
                    </NFormItem>
                    <NFormItem label="Link Name" path="name">
                        <NInput v-model:value="linkData.name" />
                    </NFormItem>
                    <NFormItem label="Link Description" path="link_desc">
                        <NInput v-model:value="linkData.link_desc" />
                    </NFormItem>
                    <NFormItem class="w-full px-7" label="Tags" path="link_tags">
                        <NSelect class="w-full" size="large" v-model:value="linkData.link_tags" multiple filterable
                            laceholder="Search Tags" :options="tagsRef" :loading="loading" clearable remote
                            :clear-filter-after-select="false" @search="handleSearch" label-field="name"
                            value-field="ID" />
                    </NFormItem>
                    <NButton type="primary" @click="validateAndSaveLink">
                        <span>Add Tag</span>
                        <template #icon>
                            <Add />
                        </template>
                    </NButton>
                </div>
            </NForm>

        </NCard>
    </NModal>
</template>