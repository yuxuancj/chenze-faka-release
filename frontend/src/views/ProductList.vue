<template>
    <Layout>
        <div class="space-y-4">
            <div class="flex flex-col md:flex-row md:items-center gap-3">
                <h2 class="text-xl font-bold text-gray-800 mr-auto">商品列表</h2>
                <div class="flex items-center gap-2">
                    <input
                        v-model="keyword"
                        type="text"
                        placeholder="搜索商品"
                        class="form-input w-48"
                        @keyup.enter="search"
                    >
                    <button @click="search" class="btn-primary">搜索</button>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else-if="products.length === 0" class="card p-8 text-center text-gray-500">
                暂无商品
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                <router-link
                    v-for="item in products"
                    :key="item.id"
                    :to="'/product/' + item.id"
                    class="card hover:shadow-md transition-shadow overflow-hidden block"
                >
                    <div class="h-48 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center text-gray-400 text-sm p-4 text-center">
                        {{ item.name }}
                    </div>
                    <div class="card-body">
                        <h3 class="font-semibold text-gray-800 truncate">{{ item.name }}</h3>
                        <div class="text-lg font-bold text-blue-600 mt-2">￥{{ item.price }}</div>
                        <div class="text-xs text-gray-500 mt-1">库存: {{ item.stock || 0 }}</div>
                    </div>
                </router-link>
            </div>

            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center gap-2 mt-6">
                <button @click="prevPage" :disabled="page <= 1" class="btn-secondary btn-sm">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-secondary btn-sm">下一页</button>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { productList } from '../api/product'

const products = ref([])
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

const totalPages = computed(() => Math.ceil(pagination.total / pagination.size) || 1)

function loadProducts() {
    loading.value = true
    productList(page.value, pageSize, keyword.value).then((data) => {
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
    loadProducts()
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        loadProducts()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        loadProducts()
    }
}

onMounted(loadProducts)
</script>
