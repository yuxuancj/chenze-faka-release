<template>
    <AdminLayout page-title="分类管理">
        <el-card shadow="never" class="mb-4">
            <div class="flex items-center justify-between">
                <span class="text-gray-700">商品分类列表</span>
                <el-button type="primary" @click="showAddModal">新增分类</el-button>
            </div>
        </el-card>

        <el-card v-loading="loading" shadow="never">
            <el-table :data="categories" style="width: 100%" empty-text="暂无分类" stripe>
                <el-table-column prop="id" label="ID" width="120" />
                <el-table-column prop="name" label="名称" min-width="200" />
                <el-table-column prop="sort" label="排序" width="120" />
                <el-table-column label="操作" width="200" fixed="right">
                    <template #default="scope">
                        <el-button size="small" type="primary" @click="editItem(scope.row)">编辑</el-button>
                        <el-button size="small" type="danger" @click="del(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>

        <el-dialog
            v-model="showModal"
            :title="editing.id ? '编辑分类' : '新增分类'"
            width="500px"
        >
            <el-form :model="editing" label-width="80px">
                <el-form-item label="分类名称">
                    <el-input v-model="editing.name" placeholder="请输入分类名称" />
                </el-form-item>
                <el-form-item label="排序">
                    <el-input-number v-model="editing.sort" :step="1" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="cancelEdit">取消</el-button>
                <el-button type="primary" :loading="saving" @click="save">保存</el-button>
            </template>
        </el-dialog>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminCategoryList, adminCategoryCreate, adminCategoryUpdate, adminCategoryDelete } from '../api/admin'

const categories = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editing = reactive({ id: null, name: '', sort: 0 })

function load() {
    loading.value = true
    adminCategoryList().then((data) => {
        if (Array.isArray(data)) {
            categories.value = data
        } else if (data && data.list) {
            categories.value = data.list
        } else {
            categories.value = []
        }
    }).catch(() => {
        categories.value = []
    }).finally(() => {
        loading.value = false
    })
}

function showAddModal() {
    editing.id = null
    editing.name = ''
    editing.sort = 0
    showModal.value = true
}

function editItem(item) {
    editing.id = item.id
    editing.name = item.name
    editing.sort = item.sort || 0
    showModal.value = true
}

function cancelEdit() {
    showModal.value = false
    editing.id = null
    editing.name = ''
    editing.sort = 0
}

function save() {
    if (!editing.name) {
        ElMessage.warning('请输入分类名称')
        return
    }
    saving.value = true
    const action = editing.id
        ? adminCategoryUpdate(editing.id, { name: editing.name, sort: editing.sort })
        : adminCategoryCreate({ name: editing.name, sort: editing.sort })
    action.then(() => {
        ElMessage.success('保存成功')
        cancelEdit()
        load()
    }).catch(() => {}).finally(() => {
        saving.value = false
    })
}

function del(row) {
    ElMessageBox.confirm('确认删除该分类？', '提示', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        adminCategoryDelete(row.id).then(() => {
            ElMessage.success('删除成功')
            load()
        }).catch(() => {})
    }).catch(() => {})
}

onMounted(load)
</script>
