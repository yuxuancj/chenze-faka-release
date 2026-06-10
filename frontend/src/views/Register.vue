<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4 register-wrapper">
        <el-card class="register-card" shadow="always">
            <template #header>
                <h2 class="text-xl font-bold text-gray-800 text-center">用户注册</h2>
            </template>
            <el-form
                ref="formRef"
                :model="form"
                :rules="rules"
                label-width="100px"
                size="large"
            >
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="form.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item label="昵称" prop="nickname">
                    <el-input v-model="form.nickname" placeholder="请输入昵称" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="form.password" type="password" show-password placeholder="请输入密码（至少6位）" />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirm">
                    <el-input v-model="form.confirm" type="password" show-password placeholder="请再次输入密码" @keyup.enter="submit" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" class="w-full" :loading="loading" @click="submit">
                        注册
                    </el-button>
                </el-form-item>
                <div class="flex items-center justify-between text-sm">
                    <router-link to="/user/login" class="text-blue-500 hover:text-blue-600">已有账号？去登录</router-link>
                    <router-link to="/" class="text-blue-500 hover:text-blue-600">返回首页</router-link>
                </div>
            </el-form>
        </el-card>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '../api/user'

const router = useRouter()
const loading = ref(false)
const formRef = ref(null)
const form = reactive({ email: '', nickname: '', password: '', confirm: '' })

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

const validateConfirm = (rule, value, callback) => {
    if (!value) {
        callback(new Error('请再次输入密码'))
    } else if (value !== form.password) {
        callback(new Error('两次密码输入不一致'))
    } else {
        callback()
    }
}

const rules = reactive({
    email: [{ validator: validateEmail, trigger: 'blur' }],
    nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码至少6位', trigger: 'blur' }
    ],
    confirm: [{ validator: validateConfirm, trigger: 'blur' }]
})

function submit() {
    formRef.value.validate((valid) => {
        if (!valid) return
        loading.value = true
        register({
            email: form.email,
            password: form.password,
            nickname: form.nickname
        }).then(() => {
            ElMessage.success('注册成功，请登录')
            router.push('/user/login')
        }).catch(() => {}).finally(() => {
            loading.value = false
        })
    })
}
</script>

<style scoped>
.register-wrapper {
    min-height: 100vh;
}
.register-card {
    width: 100%;
    max-width: 460px;
}
</style>
