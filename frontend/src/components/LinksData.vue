<script setup lang="js">
import api from '@/net/api';
import { NInfiniteScroll } from 'naive-ui';
import { reactive } from 'vue';

const dataParams = reactive({
    page: 1,
    limit: 20,
    keyword: '',
    total: 0
})

const getLinks = async () => {
    const { data, status } = await api.get('/api/links', {
        params: dataParams
    })


    if (status === 200) {
        console.log(data)
    }
}


</script>


<template>
    <NInfiniteScroll></NInfiniteScroll>
</template>