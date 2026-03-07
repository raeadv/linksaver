<script setup lang="js">
import api from '@/net/api';
import { ref } from 'vue';
import { RouterLink } from 'vue-router';
import { useAuthStore } from '@/store/authStore';
import { appUtils } from '@/utils/app.utils';
import { NCard, NButton, NFormItem, NInput, NForm } from 'naive-ui'

const store = useAuthStore()
const { message, router } = appUtils()
const loading = ref(false)
const formRef = ref(null)
const formValue = ref({
    user: {
        username: '',
        password: ''
    },
})


const LoginAttempt = async () => {
    loading.value = true
    const { data, status } = await api.post('/api/auth/login', formValue.value.user)

    if (status === 200) {
        store.setLogin(data.token)
        message.success('Login Sukses')
        setTimeout(() => {
            router.replace('/')
        }, 500);
    }

    loading.value = false
}


const rules = {
    user: {
        username: {
            required: true,
            message: 'Username tidak boleh kosong',
            trigger: 'blur'
        },
        password: {
            required: true,
            message: 'Password tidak boleh kosong',
            trigger: ['input', 'blur']
        }
    }
}


const validateLogin = (e) => {
    e.preventDefault()
    formRef.value?.validate(async (errors) => {
        if (!errors) {
            await LoginAttempt()
        }
        else {
            message.error('Gagal Login')
        }
    })
}


</script>
<template>
    <n-card title="Login" size="small" style="max-width: 23rem;">
        <n-form ref="formRef" size="large" :label-width="80" :model="formValue" :rules="rules"
            :disabled="loading === true">
            <n-form-item label="Username" path="user.username">
                <n-input v-model:value="formValue.user.username" placeholder="Username" />
            </n-form-item>
            <n-form-item label="Password" path="user.password">
                <n-input type="password" show-password-on="click" v-model:value="formValue.user.password"
                    placeholder="Password" />
            </n-form-item>
            <div class="w-full flex items-center justify-around">
                <n-form-item>
                    <RouterLink to="/register">Register</RouterLink>
                </n-form-item>
                <n-form-item>
                    <n-button type="primary" @click="validateLogin">
                        Login
                    </n-button>
                </n-form-item>
            </div>
        </n-form>
    </n-card>
</template>