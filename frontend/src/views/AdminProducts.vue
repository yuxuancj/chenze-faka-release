<template>
    <AdminLayout page-title="商品管理">
        <div class="space-y-4">
            <div class="flex flex-col md:flex-row md:items-center gap-3">
                <div class="flex items-center gap-2 flex-1">
                    <input
                        v-model="keyword"
                        type="text"
                        placeholder="搜索商品"
                        class="form-input w-48"
                        @keyup.enter="search"
                    >
                    <button @click="search" class="btn-primary btn-sm">搜索</button>
                </div>
                <router-link to="/admin/product/new" class="btn-primary">新增商品</router-link>
            </div>
            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else class="card">
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>名称</th>
                                <th>价格</th>
                                <th>库存</th>
                                <th>已售</th>
                                <th>状态</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="products.length === 0">
                                <td colspan="7" class="text-center text-gray-500 py-8">暂无数据</td>
                            </tr>
                            <tr v-for="item in products" :key="item.id">
                                <td>{{ item.id }}</td>
                                <td>{{ item.name }}</td>
                                <td>￥{{ item.price }}</td>
                                <td>{{ item.stock || 0 }}</td>
                                <td>{{ item.sales || 0 }}</td>
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
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center gap-2">
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
import { adminProductList, adminProductDelete } from '../api/admin'

const products = ref([])
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

const totalPages = computed(() => Math.ceil(pagination.total / pagination.size) || 1)

function load() {
    loading.value = true
    adminProductList(page.value, pageSize, keyword.value).then((data) => {
        products.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize
    }).catch(() => {
        products.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function search() {
    page.value = 1
    load()
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

onMounted(load)
</script>
