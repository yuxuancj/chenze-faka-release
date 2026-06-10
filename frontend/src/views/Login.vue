<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div class="card w-full max-w-md">
            <div class="card-header text-center">
                <h2 class="text-xl font-bold text-gray-800">用户登录</h2>
            </div>
            <div class="card-body space-y-4">
                <div>
                    <label class="form-label">邮箱</label>
                    <input v-model="form.email" type="email" class="form-input" placeholder="请输入邮箱">
                </div>
                <div>
                    <label class="form-label">密码</label>
                    <input v-model="form.password" type="password" class="form-input" placeholder="请输入密码" @keyup.enter="submit">
                </div>
                <button @click="submit" :disabled="loading" class="btn-primary w-full">
                    {{ loading ? '登录中...' : '登录' }}
                </button>
                <div class="flex items-center justify-between text-sm">
                    <router-link to="/user/register" class="link">还没有账号？去注册</router-link>
                    <router-link to="/" class="link">返回首页</router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { login } from '../api/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)
const form = reactive({ email: '', password: '' })

function submit() {
    if (!form.email || !form.password) {
        alert('请输入邮箱和密码')
        return
    }
    const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRe.test(form.email)) {
        alert('邮箱格式不正确')
        return
    }
    loading.value = true
    login(form.email, form.password).then((data) => {
        if (data && data.token) {
            userStore.login(data.token, data.user)
            if (data.user && data.user.is_admin) {
                localStorage.setItem('is_admin', 'true')
            } else {
                localStorage.removeItem('is_admin')
            }
            const redirect = route.query.redirect || '/'
            router.push(redirect)
        }
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}
</script>
