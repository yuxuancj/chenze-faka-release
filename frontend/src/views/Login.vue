<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4 login-wrapper">
        <el-card class="login-card" shadow="always">
            <template #header>
                <h2 class="text-xl font-bold text-gray-800 text-center">用户登录</h2>
            </template>
            <el-form
                ref="formRef"
                :model="form"
                :rules="rules"
                label-width="80px"
                size="large"
            >
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="form.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" @keyup.enter="submit" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" class="w-full" :loading="loading" @click="submit">
                        登录
                    </el-button>
                </el-form-item>
                <div class="flex items-center justify-between text-sm">
                    <router-link to="/user/register" class="text-blue-500 hover:text-blue-600">还没有账号？去注册</router-link>
                    <router-link to="/" class="text-blue-500 hover:text-blue-600">返回首页</router-link>
                </div>
            </el-form>
        </el-card>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { login } from '../api/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)
const formRef = ref(null)
const form = reactive({ email: '', password: '' })

const validateEmail = (rule, value, callback) => {
    const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!value) {
        callback(new Error('请输入邮箱'))
    } else if (!emailRe.test(value)) {
        callback(new Error('邮箱格式不正确'))
    } else {
        callback()
    }
}

const rules = reactive({
    email: [{ validator: validateEmail, trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

function submit() {
    formRef.value.validate((valid) => {
        if (!valid) return
        loading.value = true
        login(form.email, form.password).then((data) => {
            if (data && data.token) {
                userStore.login(data.token, data.user)
                if (data.user && data.user.is_admin) {
                    localStorage.setItem('is_admin', 'true')
                } else {
                    localStorage.removeItem('is_admin')
                }
                ElMessage.success('登录成功')
                const redirect = route.query.redirect || '/'
                router.push(redirect)
            }
        }).catch(() => {}).finally(() => {
            loading.value = false
        })
    })
}
</script>

<style scoped>
.login-wrapper {
    min-height: 100vh;
}
.login-card {
    width: 100%;
    max-width: 420px;
}
</style>
