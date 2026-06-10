<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">个人中心</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="card">
                    <div class="card-header font-semibold">基本信息</div>
                    <div class="card-body space-y-4">
                        <div>
                            <label class="form-label">邮箱</label>
                            <input :value="user.email" type="email" class="form-input bg-gray-50" disabled>
                        </div>
                        <div>
                            <label class="form-label">昵称</label>
                            <input v-model="editNickname" type="text" class="form-input">
                        </div>
                        <button @click="saveProfile" :disabled="loading" class="btn-primary">
                            {{ loading ? '保存中...' : '保存' }}
                        </button>
                    </div>
                </div>
                <div class="card">
                    <div class="card-header font-semibold">账号信息</div>
                    <div class="card-body">
                        <table class="table">
                            <tbody>
                                <tr><td class="text-gray-500 w-24">用户ID</td><td>{{ user.id || '-' }}</td></tr>
                                <tr><td class="text-gray-500">余额</td><td>￥{{ user.balance || 0 }}</td></tr>
                                <tr><td class="text-gray-500">积分</td><td>{{ user.points || 0 }}</td></tr>
                                <tr><td class="text-gray-500">等级</td><td>{{ user.level || '-' }}</td></tr>
                                <tr><td class="text-gray-500">状态</td><td>{{ user.status || '正常' }}</td></tr>
                                <tr><td class="text-gray-500">注册时间</td><td>{{ user.created_at || '-' }}</td></tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="card md:col-span-2">
                    <div class="card-header font-semibold">修改密码</div>
                    <div class="card-body space-y-4">
                        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                            <div>
                                <label class="form-label">当前密码</label>
                                <input v-model="pwdForm.old" type="password" class="form-input">
                            </div>
                            <div>
                                <label class="form-label">新密码</label>
                                <input v-model="pwdForm.new" type="password" class="form-input">
                            </div>
                            <div>
                                <label class="form-label">确认新密码</label>
                                <input v-model="pwdForm.confirm" type="password" class="form-input" @keyup.enter="changePwd">
                            </div>
                        </div>
                        <button @click="changePwd" :disabled="pwdLoading" class="btn-primary">
                            {{ pwdLoading ? '提交中...' : '修改密码' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { profile, updateProfile, changePassword } from '../api/user'

const user = ref({})
const editNickname = ref('')
const loading = ref(false)
const pwdLoading = ref(false)
const pwdForm = reactive({ old: '', new: '', confirm: '' })

function loadProfile() {
    profile().then((res) => {
        user.value = res.data || {}
        editNickname.value = user.value.nickname || ''
    }).catch(() => {})
}

function saveProfile() {
    if (!editNickname.value) {
        alert('请输入昵称')
        return
    }
    loading.value = true
    updateProfile(editNickname.value).then(() => {
        alert('保存成功')
        loadProfile()
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

function changePwd() {
    if (!pwdForm.old || !pwdForm.new || !pwdForm.confirm) {
        alert('请完整填写密码')
        return
    }
    if (pwdForm.new !== pwdForm.confirm) {
        alert('两次输入的新密码不一致')
        return
    }
    pwdLoading.value = true
    changePassword(pwdForm.old, pwdForm.new).then(() => {
        alert('密码修改成功')
        pwdForm.old = ''
        pwdForm.new = ''
        pwdForm.confirm = ''
    }).catch(() => {}).finally(() => {
        pwdLoading.value = false
    })
}

onMounted(loadProfile)
</script>
