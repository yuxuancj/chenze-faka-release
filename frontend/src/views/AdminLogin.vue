<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div class="card w-full max-w-md">
            <div class="card-header text-center">
                <h2 class="text-xl font-bold text-gray-800">管理员登录</h2>
            </div>
            <div class="card-body space-y-4">
                <div>
                    <label class="form-label">邮箱</label>
                    <input v-model="form.email" type="email" class="form-input" placeholder="请输入管理员邮箱">
                </div>
                <div>
                    <label class="form-label">密码</label>
                    <input v-model="form.password" type="password" class="form-input" placeholder="请输入密码" @keyup.enter="submit">
                </div>
                <button @click="submit" :disabled="loading" class="btn-primary w-full">
                    {{ loading ? '登录中...' : '登录' }}
                </button>
                <div class="text-center text-sm">
                    <router-link to="/" class="link">返回首页</router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { login } from '../api/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const form = reactive({ email: '', password: '' })

function submit() {
    if (!form.email || !form.password) {
        alert('请输入邮箱和密码')
        return
    }
    loading.value = true
    login(form.email, form.password).then((data) => {
        if (data && data.token && data.user && data.user.is_admin) {
            userStore.login(data.token, data.user)
            localStorage.setItem('is_admin', 'true')
            router.push('/admin/')
        } else {
            alert('该账号没有管理员权限')
        }
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}
</script>
