<template>
    <AdminLayout page-title="分类管理">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <h3 class="font-semibold text-gray-800">商品分类</h3>
                <button @click="showAdd = true" class="btn-primary">新增分类</button>
            </div>
            <div class="card">
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>名称</th>
                                <th>排序</th>
                                <th>创建时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="categories.length === 0">
                                <td colspan="5" class="text-center text-gray-500 py-8">暂无分类</td>
                            </tr>
                            <tr v-for="item in categories" :key="item.id">
                                <td>{{ item.id }}</td>
                                <td>{{ item.name }}</td>
                                <td>{{ item.sort || 0 }}</td>
                                <td>{{ item.created_at }}</td>
                                <td>
                                    <button @click="editItem(item)" class="btn-sm btn-primary">编辑</button>
                                    <button @click="del(item.id)" class="btn-sm btn-danger ml-2">删除</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div v-if="showAdd" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20">
                <div class="card w-full max-w-md">
                    <div class="card-header font-semibold">{{ editing.id ? '编辑分类' : '新增分类' }}</div>
                    <div class="card-body space-y-4">
                        <div>
                            <label class="form-label">分类名称</label>
                            <input v-model="editing.name" type="text" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">排序</label>
                            <input v-model.number="editing.sort" type="number" class="form-input">
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
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminCategoryList, adminCategoryCreate, adminCategoryUpdate, adminCategoryDelete } from '../api/admin'

const categories = ref([])
const showAdd = ref(false)
const loading = ref(false)
const editing = reactive({ id: null, name: '', sort: 0 })

function load() {
    adminCategoryList(1, 100).then((res) => {
        categories.value = res.data && res.data.list ? res.data.list : []
    }).catch(() => {})
}

function editItem(item) {
    editing.id = item.id
    editing.name = item.name
    editing.sort = item.sort || 0
    showAdd.value = true
}

function cancel() {
    showAdd.value = false
    editing.id = null
    editing.name = ''
    editing.sort = 0
}

function save() {
    if (!editing.name) {
        alert('请输入分类名称')
        return
    }
    loading.value = true
    const action = editing.id
        ? adminCategoryUpdate(editing.id, { name: editing.name, sort: editing.sort })
        : adminCategoryCreate({ name: editing.name, sort: editing.sort })
    action.then(() => {
        alert('保存成功')
        cancel()
        load()
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

function del(id) {
    if (!confirm('确认删除该分类？')) return
    adminCategoryDelete(id).then(() => {
        load()
    }).catch(() => {})
}

onMounted(load)
</script>
