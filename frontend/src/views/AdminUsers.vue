<template>
    <AdminLayout page-title="用户管理">
        <div class="space-y-4">
            <div class="card">
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>邮箱</th>
                                <th>昵称</th>
                                <th>余额</th>
                                <th>积分</th>
                                <th>等级</th>
                                <th>状态</th>
                                <th>注册时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="users.length === 0">
                                <td colspan="9" class="text-center text-gray-500 py-8">暂无用户</td>
                            </tr>
                            <tr v-for="user in users" :key="user.id">
                                <td>{{ user.id }}</td>
                                <td>{{ user.email }}</td>
                                <td>{{ user.nickname }}</td>
                                <td>￥{{ user.balance || 0 }}</td>
                                <td>{{ user.points || 0 }}</td>
                                <td>{{ user.level || 0 }}</td>
                                <td>
                                    <span v-if="user.status === 1" class="badge-green">正常</span>
                                    <span v-else class="badge-red">禁用</span>
                                </td>
                                <td>{{ user.created_at }}</td>
                                <td>
                                    <button @click="editUser(user)" class="btn-sm btn-primary">编辑</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center space-x-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>

            <div v-if="editing.id" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20">
                <div class="card w-full max-w-md">
                    <div class="card-header font-semibold">编辑用户</div>
                    <div class="card-body space-y-4">
                        <div>
                            <label class="form-label">昵称</label>
                            <input v-model="editing.nickname" type="text" class="form-input">
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label class="form-label">余额</label>
                                <input v-model.number="editing.balance" type="number" step="0.01" class="form-input">
                            </div>
                            <div>
                                <label class="form-label">积分</label>
                                <input v-model.number="editing.points" type="number" class="form-input">
                            </div>
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label class="form-label">等级</label>
                                <input v-model.number="editing.level" type="number" class="form-input">
                            </div>
                            <div>
                                <label class="form-label">状态</label>
                                <select v-model="editing.status" class="form-input">
                                    <option :value="1">正常</option>
                                    <option :value="0">禁用</option>
                                </select>
                            </div>
                        </div>
                        <div class="flex items-center justify-end space-x-2">
                            <button @click="cancel" class="btn-secondary">取消</button>
                            <button @click="save" :disabled="loading" class="btn-primary">
                                {{ loading ? '保存中...' : '保存' }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminUserList, adminUserUpdate } from '../api/admin'

const users = ref([])
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })
const loading = ref(false)
const editing = reactive({ id: null, nickname: '', balance: 0, points: 0, level: 0, status: 1 })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function load() {
    adminUserList(page.value, pageSize).then((res) => {
        users.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
}

function editUser(user) {
    editing.id = user.id
    editing.nickname = user.nickname || ''
    editing.balance = user.balance || 0
    editing.points = user.points || 0
    editing.level = user.level || 0
    editing.status = user.status || 1
}

function cancel() {
    editing.id = null
    editing.nickname = ''
    editing.balance = 0
    editing.points = 0
    editing.level = 0
    editing.status = 1
}

function save() {
    loading.value = true
    adminUserUpdate(editing.id, {
        nickname: editing.nickname,
        balance: editing.balance,
        points: editing.points,
        level: editing.level,
        status: editing.status
    }).then(() => {
        alert('保存成功')
        cancel()
        load()
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        load()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        load()
    }
}

onMounted(load)
</script>
