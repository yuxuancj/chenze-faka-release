<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">个人中心</h2>

            <el-row v-loading="loading" :gutter="20">
                <el-col :xs="24" :md="12">
                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">基本信息</span>
                        </template>
                        <el-form
                            ref="profileFormRef"
                            :model="profileForm"
                            :rules="profileRules"
                            label-width="80px"
                        >
                            <el-form-item label="邮箱">
                                <el-input v-model="profileForm.email" disabled />
                            </el-form-item>
                            <el-form-item label="昵称" prop="nickname">
                                <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" :loading="saving" @click="saveProfile">
                                    保存
                                </el-button>
                            </el-form-item>
                        </el-form>
                    </el-card>
                </el-col>

                <el-col :xs="24" :md="12">
                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">账号信息</span>
                        </template>
                        <el-table :data="accountInfo" border style="width: 100%">
                            <el-table-column prop="label" label="项目" width="120" />
                            <el-table-column prop="value" label="内容" />
                        </el-table>
                    </el-card>
                </el-col>

                <el-col :span="24">
                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">修改密码</span>
                        </template>
                        <el-form
                            ref="pwdFormRef"
                            :model="pwdForm"
                            :rules="pwdRules"
                            label-width="120px"
                        >
                            <el-row :gutter="20">
                                <el-col :xs="24" :sm="8">
                                    <el-form-item label="当前密码" prop="old">
                                        <el-input v-model="pwdForm.old" type="password" show-password placeholder="请输入当前密码" />
                                    </el-form-item>
                                </el-col>
                                <el-col :xs="24" :sm="8">
                                    <el-form-item label="新密码" prop="new">
                                        <el-input v-model="pwdForm.new" type="password" show-password placeholder="请输入新密码" />
                                    </el-form-item>
                                </el-col>
                                <el-col :xs="24" :sm="8">
                                    <el-form-item label="确认新密码" prop="confirm">
                                        <el-input v-model="pwdForm.confirm" type="password" show-password placeholder="请再次输入新密码" @keyup.enter="changePwd" />
                                    </el-form-item>
                                </el-col>
                            </el-row>
                            <el-form-item>
                                <el-button type="primary" :loading="pwdLoading" @click="changePwd">
                                    修改密码
                                </el-button>
                            </el-form-item>
                        </el-form>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Layout from '../components/Layout.vue'
import { profile, updateProfile, changePassword } from '../api/user'

const user = ref({})
const loading = ref(true)
const saving = ref(false)
const pwdLoading = ref(false)
const profileFormRef = ref(null)
const pwdFormRef = ref(null)

const profileForm = reactive({ email: '', nickname: '' })
const pwdForm = reactive({ old: '', new: '', confirm: '' })

const profileRules = reactive({
    nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }]
})

const validateConfirm = (rule, value, callback) => {
    if (!value) {
        callback(new Error('请再次输入新密码'))
    } else if (value !== pwdForm.new) {
        callback(new Error('两次输入的新密码不一致'))
    } else {
        callback()
    }
}

const pwdRules = reactive({
    old: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
    new: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '新密码至少6位', trigger: 'blur' }
    ],
    confirm: [{ validator: validateConfirm, trigger: 'blur' }]
})

const accountInfo = computed(() => [
    { label: '用户ID', value: user.value.id || '-' },
    { label: '余额', value: '￥' + (user.value.balance || 0) },
    { label: '积分', value: user.value.points || 0 },
    { label: '等级', value: user.value.level || '-' },
    { label: '状态', value: user.value.status ? '正常' : '禁用' },
    { label: '注册时间', value: user.value.created_at || '-' }
])

function loadProfile() {
    loading.value = true
    profile().then((data) => {
        user.value = data || {}
        profileForm.email = user.value.email || ''
        profileForm.nickname = user.value.nickname || ''
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

function saveProfile() {
    profileFormRef.value.validate((valid) => {
        if (!valid) return
        saving.value = true
        updateProfile(profileForm.nickname).then(() => {
            ElMessage.success('保存成功')
            loadProfile()
        }).catch(() => {}).finally(() => {
            saving.value = false
        })
    })
}

function changePwd() {
    pwdFormRef.value.validate((valid) => {
        if (!valid) return
        pwdLoading.value = true
        changePassword(pwdForm.old, pwdForm.new).then(() => {
            ElMessage.success('密码修改成功')
            pwdForm.old = ''
            pwdForm.new = ''
            pwdForm.confirm = ''
            pwdFormRef.value.resetFields()
        }).catch(() => {}).finally(() => {
            pwdLoading.value = false
        })
    })
}

onMounted(loadProfile)
</script>
