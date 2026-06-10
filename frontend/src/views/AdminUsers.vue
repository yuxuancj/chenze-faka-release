<template>
    <AdminLayout page-title="用户管理">
        <el-table v-loading="loading" :data="users" style="width: 100%" border stripe empty-text="暂无用户">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column prop="nickname" label="昵称" />
            <el-table-column label="余额" width="120">
                <template #default="scope">￥{{ scope.row.balance || 0 }}</template>
            </el-table-column>
            <el-table-column prop="points" label="积分" width="100" />
            <el-table-column prop="level" label="等级" width="80" />
            <el-table-column label="状态" width="100">
                <template #default="scope">
                    <el-tag v-if="scope.row.status === 1 || scope.row.status === undefined" type="success">正常</el-tag>
                    <el-tag v-else type="danger">禁用</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="created_at" label="注册时间" width="180" />
            <el-table-column label="操作" width="120" fixed="right">
                <template #default="scope">
                    <el-button type="primary" size="small" @click="editUser(scope.row)">编辑</el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-pagination
            v-if="pagination.total > pagination.size"
            class="mt-4"
            v-model:current-page="page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            background
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />

        <el-dialog v-model="dialogVisible" title="编辑用户" width="500px">
            <el-form :model="editing" label-width="100px">
                <el-form-item label="昵称">
                    <el-input v-model="editing.nickname" />
                </el-form-item>
                <el-form-item label="余额">
                    <el-input-number v-model="editing.balance" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="积分">
                    <el-input-number v-model="editing.points" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="等级">
                    <el-input-number v-model="editing.level" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="editing.status" style="width: 100%">
                        <el-option label="正常" :value="1" />
                        <el-option label="禁用" :value="0" />
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="cancel">取消</el-button>
                <el-button type="primary" :loading="saving" @click="save">保存</el-button>
            </template>
        </el-dialog>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminUserList, adminUserUpdate } from '../api/admin'

const users = ref([])
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editing = reactive({ id: null, nickname: '', balance: 0, points: 0, level: 0, status: 1 })

function load() {
    loading.value = true
    adminUserList(page.value, pagination.size).then((data) => {
        users.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize
    }).catch(() => {
        users.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function editUser(user) {
    editing.id = user.id
    editing.nickname = user.nickname || ''
    editing.balance = user.balance || 0
    editing.points = user.points || 0
    editing.level = user.level || 1
    editing.status = user.status === undefined ? 1 : user.status
    dialogVisible.value = true
}

function cancel() {
    dialogVisible.value = false
    editing.id = null
    editing.nickname = ''
    editing.balance = 0
    editing.points = 0
    editing.level = 0
    editing.status = 1
}

function save() {
    saving.value = true
    adminUserUpdate(editing.id, {
        nickname: editing.nickname,
        balance: editing.balance,
        points: editing.points,
        level: editing.level,
        status: editing.status
    }).then(() => {
        ElMessage.success('保存成功')
        cancel()
        load()
    }).catch(() => {}).finally(() => {
        saving.value = false
    })
}

function handleSizeChange(val) {
    pagination.size = val
    page.value = 1
    load()
}

function handleCurrentChange(val) {
    page.value = val
    load()
}

onMounted(load)
</script>
