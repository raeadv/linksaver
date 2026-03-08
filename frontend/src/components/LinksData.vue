<script setup lang="js">
import api from '@/net/api';
import { NInfiniteScroll, NSpin } from 'naive-ui';
import { onMounted, ref, watch } from 'vue';


const LIMIT = 20

const links = ref([])
const offset = ref(0)
const keyword = ref('')
const hasMore = ref(true)

const fetchLinks = async (done) => {
    try {
        const { data, status } = await api.get('/api/links/scroll', {
            params: {
                offset: offset.value,
                limit: LIMIT,
                keyword: keyword.value,
            }
        })

        if (status === 200 && data.status) {
            console.log(data)
            links.value.push(...data.links)
            offset.value += data.links.length
            hasMore.value = data.has_more
        } else {
            hasMore.value = false
        }
    } catch {
        hasMore.value = false
    } finally {
        done(!hasMore.value)
    }
}

const reset = async () => {
    links.value = []
    offset.value = 0
    hasMore.value = true
    await fetchLinks(() => { })
}

defineExpose({ reset, fetchLinks })

watch(keyword, () => reset())
onMounted(() => fetchLinks(() => { }))
</script>

<template>
    <div class="w-full">
        <input v-model="keyword" type="text" placeholder="Search links..."
            class="w-full mb-3 px-3 py-2 border rounded" />
        <NInfiniteScroll :distance="100" @load="fetchLinks">
            <div v-for="link in links" :key="link.ID" class="p-3 mb-2 border rounded">
                <a :href="link.link" target="_blank" class="font-medium text-blue-600 hover:underline">
                    {{ link.name || link.link }}
                </a>
                <p v-if="link.link_desc" class="text-sm text-gray-500 mt-1">{{ link.link_desc }}</p>
                <div v-if="link.LinkTags && link.LinkTags.length" class="flex flex-wrap gap-1 mt-1">
                    <span v-for="lt in link.LinkTags" :key="lt.TagId" class="text-xs bg-gray-100 px-2 py-0.5 rounded">{{
                        lt.Tag?.name }}</span>
                </div>
            </div>

            <template #loading>
                <div class="flex justify-center py-4">
                    <NSpin size="small" />
                </div>
            </template>
        </NInfiniteScroll>

        <div v-if="!hasMore && links.length === 0" class="text-center text-gray-400 py-8">
            No links found.
        </div>
    </div>
</template>
