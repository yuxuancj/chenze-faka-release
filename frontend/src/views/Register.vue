<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div class="card w-full max-w-md">
            <div class="card-header text-center">
                <h2 class="text-xl font-bold text-gray-800">用户注册</h2>
            </div>
            <div class="card-body space-y-4">
                <div>
                    <label class="form-label">邮箱</label>
                    <input v-model="form.email" type="email" class="form-input" placeholder="请输入邮箱">
                </div>
                <div>
                    <label class="form-label">昵称</label>
                    <input v-model="form.nickname" type="text" class="form-input" placeholder="请输入昵称">
                </div>
                <div>
                    <label class="form-label">密码</label>
                    <input v-model="form.password" type="password" class="form-input" placeholder="请输入密码">
                </div>
                <div>
                    <label class="form-label">确认密码</label>
                    <input v-model="form.confirm" type="password" class="form-input" placeholder="请再次输入密码" @keyup.enter="submit">
                </div>
                <button @click="submit" :disabled="loading" class="btn-primary w-full">
                    {{ loading ? '注册中...' : '注册' }}
                </button>
                <div class="flex items-center justify-between text-sm">
                    <router-link to="/user/login" class="link">已有账号？去登录</router-link>
                    <router-link to="/" class="link">返回首页</router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/user'

const router = useRouter()
const loading = ref(false)
const form = reactive({ email: '', nickname: '', password: '', confirm: '' })

function submit() {
    if (!form.email || !form.password || !form.nickname) {
        alert('请填写完整信息')
        return
    }
    if (form.password !== form.confirm) {
        alert('两次密码输入不一致')
        return
    }
    loading.value = true
    register(form.email, form.password, form.nickname).then((res) => {
        alert('注册成功，请登录')
        router.push('/user/login')
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}
</script>
