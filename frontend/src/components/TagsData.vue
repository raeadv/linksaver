<script setup lang="js">
import api from '@/net/api';
import { NFormItem, NInput, NPagination, NSelect, NSpace, NTag } from 'naive-ui';
import { onMounted, reactive, ref, watch } from 'vue';
const debounce = (fn, delay) => { let t; return (...args) => { clearTimeout(t); t = setTimeout(() => fn(...args), delay) } }


const tagParameter = reactive({
    page: 1,
    limit: 25,
    keyword: '',
    total: 0
})

const tagsData = ref([])

const getTags = async () => {
    const { data, status } = await api.get('/api/tags', {
        params: tagParameter
    })

    if (status === 200) {
        console.log({ data, status })
        tagsData.value = data.tags
        tagParameter.total = data.metadata.total
    }


}

onMounted(() => getTags())

const debouncedGetTags = debounce(getTags, 500)

watch([() => tagParameter.limit, () => tagParameter.page], () => {
    getTags()
})

watch(() => tagParameter.keyword, (val) => {
    if (val.length === 0 || val.length > 2) {
        debouncedGetTags()
    }
})


async function handlerDeleteTag(tagId) {
    console.log("delete this : ", tagId)
}


</script>

<template>
    <div class="my-2 flex flex-row space-x-2">
        <NFormItem label="Keyword">
            <NInput v-model:value="tagParameter.keyword" placeholder="search tags" />
        </NFormItem>
        <NFormItem>
            <NSelect :options="[{ label: '25', value: 25 }, { label: '50', value: 50 }, { label: '100', value: 100 }]"
                v-model:value="tagParameter.limit" />
        </NFormItem>

    </div>
    <div class="flex flex-row space-x-2" v-if="tagsData.length > 0">
        <NSpace v-for="tag in tagsData">
            <NTag closable @close="() => handlerDeleteTag(tag.ID)" size="large" round>
                {{ tag.name }}
            </NTag>
        </NSpace>
    </div>
    <div class="my-2">
        <NPagination v-model:page="tagParameter.page" :page-count="tagParameter.total"
            :page-size="tagParameter.limit" />
    </div>

</template>