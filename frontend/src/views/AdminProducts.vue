<template>
    <AdminLayout page-title="商品管理">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <select v-model="filters.category_id" class="form-input w-32">
                        <option value="">全部分类</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                    </select>
                    <input v-model="filters.keyword" type="text" placeholder="搜索商品" class="form-input w-48" @keyup.enter="load">
                    <button @click="load" class="btn-primary btn-sm">搜索</button>
                </div>
                <router-link to="/admin/product/new" class="btn-primary">新增商品</router-link>
            </div>
            <div class="card">
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>名称</th>
                                <th>价格</th>
                                <th>库存</th>
                                <th>已售</th>
                                <th>分类</th>
                                <th>状态</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="products.length === 0">
                                <td colspan="8" class="text-center text-gray-500 py-8">暂无数据</td>
                            </tr>
                            <tr v-for="item in products" :key="item.id">
                                <td>{{ item.id }}</td>
                                <td>{{ item.name }}</td>
                                <td>￥{{ item.price }}</td>
                                <td>{{ item.stock || 0 }}</td>
                                <td>{{ item.sold || 0 }}</td>
                                <td>{{ item.category_name || '-' }}</td>
                                <td>
                                    <span v-if="item.status === 1" class="badge-green">上架</span>
                                    <span v-else class="badge-gray">下架</span>
                                </td>
                                <td>
                                    <router-link :to="'/admin/product/' + item.id" class="btn-sm btn-primary">编辑</router-link>
                                    <button @click="del(item.id)" class="btn-sm btn-danger ml-2">删除</button>
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
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminProductList, adminProductDelete, adminCategoryList } from '../api/admin'

const products = ref([])
const categories = ref([])
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })
const filters = reactive({ category_id: '', keyword: '' })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function loadCategories() {
    adminCategoryList(1, 100).then((res) => {
        categories.value = res.data && res.data.list ? res.data.list : []
    }).catch(() => {})
}

function load() {
    adminProductList(page.value, pageSize, filters.category_id || 0, filters.keyword || '').then((res) => {
        products.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
}

function del(id) {
    if (!confirm('确认删除该商品？')) return
    adminProductDelete(id).then(() => {
        load()
    }).catch(() => {})
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

onMounted(() => {
    loadCategories()
    load()
})
</script>
