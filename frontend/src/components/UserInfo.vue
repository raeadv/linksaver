<script setup lang="js">
import { useAuthStore } from '@/store/authStore';
import { NCard } from 'naive-ui';
import { computed } from 'vue';

const store = useAuthStore()
const user = computed(() => store.userdata)
</script>



<template>
    <NCard title="User Info">
        <table>
            <tr>
                <th align="left">Username</th>
                <th>:</th>
                <td>{{ user.username }}</td>
            </tr>
            <tr>
                <th align="left">Email</th>
                <th>:</th>
                <td>{{ user.email }}</td>
            </tr>
        </table>
    </NCard>

</template>