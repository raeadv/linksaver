<script setup lang="js">
import api from '@/net/api';
import { reactive, ref } from 'vue';
import { RouterLink } from 'vue-router';
import { z } from 'zod';
import { createFormRules } from '@/utils/validation.adapter';
import { useAuthStore } from '@/store/authStore';
import { appUtils } from '@/utils/app.utils';
import { NCard, NButton, NFormItem, NInput, NForm } from 'naive-ui'


const store = useAuthStore()
const loading = ref(false)
const { router, message } = appUtils()
const formRef = ref()
const formValue = reactive({
    name: '',
    username: '',
    email: '',
    phone: '',
    password: '',
    password_confirmation: '',
})


const LoginAttempt = async () => {
    loading.value = true
    const { data, status } = await api.post('/api/auth/register', formValue)

    if (status === 201) {
        store.setLogin(data)
        router.replace('/')
    }

    loading.value = false
}


const schema = z.object({
    name: z.string().min(5, { message: 'name min 3 chars' }),
    username: z.string().min(5, { message: 'Username min 5 chars' }).uppercase({ message: 'username must use all capital' }),
    email: z.email({ message: 'Invalid email format' }),
    phone: z.string().min(10, { message: 'Phone Number min 10 chars' }),
    password: z.string().min(8, { message: 'password min 8 chars' }),
    password_confirmation: z.string({ message: 'Password Confirmation must not empty' }).refine(data => data.password === data.password_confirmation, {
        message: 'Password Confirmation incorrect',
        path: ['password_confirmation']
    }),
})


const rules = createFormRules(schema)

const validateRegister = (e) => {
    e.preventDefault()
    formRef.value?.validate(async (errors) => {
        if (!errors) {
            await LoginAttempt()
        }
        else {
            message.error('Register Failed')
        }
    })
}


</script>
<template>
    <n-card title="Register" size="small" style="max-width: 23rem;">
        <n-form ref="formRef" size="large" :label-width="80" :model="formValue" :rules="rules" :disabled="loading">
            <n-form-item label="Name" path="name">
                <n-input v-model:value="formValue.name" placeholder="Name" />
            </n-form-item>
            <n-form-item label="Username" path="username">
                <n-input v-model:value="formValue.username" placeholder="Username" />
            </n-form-item>
            <n-form-item label="Email" path="email">
                <n-input v-model:value="formValue.email" placeholder="Email" />
            </n-form-item>
            <n-form-item label="Phone" path="phone">
                <n-input v-model:value="formValue.phone" placeholder="Phone" />
            </n-form-item>
            <n-form-item label="Password" path="password">
                <n-input type="password" show-password-on="click" v-model:value="formValue.password"
                    placeholder="Password" />
            </n-form-item>

            <n-form-item label="Konfirmasi Password" path="password_confirmation">
                <n-input type="password" show-password-on="click" v-model:value="formValue.password_confirmation"
                    placeholder="Password" />
            </n-form-item>
            <div class="w-full flex items-center justify-around">
                <n-form-item>
                    <RouterLink to="/login">Login</RouterLink>
                </n-form-item>
                <n-form-item>
                    <n-button type="primary" @click="validateRegister">
                        Register
                    </n-button>
                </n-form-item>
            </div>
        </n-form>
    </n-card>
</template>