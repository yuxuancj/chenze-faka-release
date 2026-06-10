<template>
    <Layout>
        <el-card class="mb-4" shadow="never">
            <div class="flex items-center justify-between flex-wrap gap-3">
                <h2 class="text-xl font-bold text-gray-800">商品列表</h2>
                <div class="flex items-center gap-2">
                    <el-input
                        v-model="keyword"
                        placeholder="搜索商品"
                        style="width: 240px"
                        clearable
                        @keyup.enter="search"
                    >
                        <template #append>
                            <el-button @click="search">搜索</el-button>
                        </template>
                    </el-input>
                </div>
            </div>
        </el-card>

        <div v-if="loading" class="text-center py-16">
            <el-icon class="is-loading text-blue-600 text-4xl"><Loading /></el-icon>
            <p class="text-gray-500 mt-4">加载中...</p>
        </div>
        <div v-else-if="products.length === 0" class="text-center py-16 text-gray-500">
            暂无商品
        </div>
        <el-row v-else :gutter="16">
            <el-col
                v-for="item in products"
                :key="item.id"
                :xs="12"
                :sm="12"
                :md="8"
                :lg="6"
                :xl="6"
            >
                <router-link :to="'/product/' + item.id">
                    <el-card class="mb-4 product-card hover-shadow cursor-pointer" shadow="hover">
                        <div class="product-placeholder flex items-center justify-center text-gray-400 text-sm p-4 text-center">
                            {{ item.name }}
                        </div>
                        <div class="p-2">
                            <h3 class="font-semibold text-gray-800 truncate mb-1">{{ item.name }}</h3>
                            <div class="flex items-baseline justify-between">
                                <span class="text-lg font-bold text-blue-600">￥{{ item.price }}</span>
                                <span class="text-xs text-gray-500">库存: {{ item.stock || 0 }}</span>
                            </div>
                        </div>
                    </el-card>
                </router-link>
            </el-col>
        </el-row>

        <div class="flex justify-center mt-6" v-if="pagination && pagination.total > 0">
            <el-pagination
                v-model:current-page="page"
                v-model:page-size="pagination.size"
                :page-sizes="[10, 20, 50, 100]"
                :total="pagination.total"
                layout="prev, pager, next, total"
                background
                @current-change="loadProducts"
                @size-change="loadProducts"
            />
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import Layout from '../components/Layout.vue'
import { productList } from '../api/product'

const products = ref([])
const page = ref(1)
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

function loadProducts() {
    loading.value = true
    productList(page.value, pagination.size, keyword.value).then((data) => {
        products.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pagination.size
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

onMounted(loadProducts)
</script>

<style scoped>
.product-card {
    transition: transform 0.2s, box-shadow 0.2s;
}
.product-card:hover {
    transform: translateY(-2px);
}
.product-placeholder {
    height: 180px;
    background: linear-gradient(135deg, #f7f7f7 0%, #e8e8e8 100%);
    border-radius: 4px;
    margin: -15px -15px 0 -15px;
}
</style>
