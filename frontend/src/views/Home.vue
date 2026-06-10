<template>
    <Layout>
        <div class="space-y-8">
            <div class="relative bg-gradient-to-r from-blue-600 to-blue-500 rounded-lg overflow-hidden">
                <div class="px-6 py-12 md:px-12 md:py-16">
                    <h1 class="text-2xl md:text-4xl font-bold text-white mb-3">商品精选</h1>
                    <p class="text-blue-100 mb-6 max-w-lg">浏览各类优质商品，快速下单，即刻到账。</p>
                    <div class="flex flex-wrap gap-3">
                        <router-link to="/products" class="bg-white text-blue-600 px-6 py-2 rounded-md font-medium hover:bg-blue-50 transition-colors">浏览全部商品</router-link>
                        <router-link to="/user/login" class="border border-white text-white px-6 py-2 rounded-md font-medium hover:bg-white hover:bg-opacity-10 transition-colors">登录账户</router-link>
                    </div>
                </div>
                <div class="absolute right-0 top-0 w-64 h-64 bg-white bg-opacity-10 rounded-full -mr-32 -mt-32"></div>
                <div class="absolute right-12 bottom-6 hidden md:block text-right">
                    <div class="text-white text-opacity-80 text-xs">优质商品</div>
                    <div class="text-white text-2xl font-bold">快速发货</div>
                </div>
            </div>

            <div class="space-y-4">
                <div class="flex items-center justify-between">
                    <h2 class="text-xl font-bold text-gray-800">热门商品</h2>
                    <router-link to="/products" class="text-sm text-blue-600 hover:underline">查看全部</router-link>
                </div>
                <div v-if="loading" class="card p-8 text-center text-gray-500">
                    加载中...
                </div>
                <div v-else-if="products.length === 0" class="card p-8 text-center text-gray-500">
                    暂无商品
                </div>
                <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <router-link
                        v-for="item in products"
                        :key="item.id"
                        :to="'/product/' + item.id"
                        class="card hover:shadow-md transition-shadow overflow-hidden block"
                    >
                        <div class="h-40 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center text-gray-400 text-sm">
                            {{ item.name }}
                        </div>
                        <div class="card-body">
                            <h3 class="font-semibold text-gray-800 truncate">{{ item.name }}</h3>
                            <div class="text-xl font-bold text-blue-600 mt-2">￥{{ item.price }}</div>
                            <div class="text-xs text-gray-500 mt-1">库存: {{ item.stock || 0 }}</div>
                        </div>
                    </router-link>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { productList } from '../api/product'

const products = ref([])
const loading = ref(true)

function load() {
    loading.value = true
    productList(1, 6, '').then((data) => {
        products.value = (data && data.list) ? data.list : []
    }).catch(() => {
        products.value = []
    }).finally(() => {
        loading.value = false
    })
}

onMounted(load)
</script>
