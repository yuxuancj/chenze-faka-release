<template>
    <Layout>
        <div class="space-y-6">
            <div class="card">
                <div class="card-body">
                    <h2 class="text-2xl font-bold text-gray-800 mb-2">欢迎来到商城</h2>
                    <p class="text-gray-600">浏览精选商品，享受便捷的购物体验。</p>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div v-for="item in [1,2,3,4,5,6]" :key="item" class="card p-4 text-center">
                    <div class="bg-gray-100 h-32 rounded-md mb-3 flex items-center justify-center text-gray-400">商品图片</div>
                    <h3 class="font-semibold text-gray-800">商品 {{ item }}</h3>
                    <p class="text-blue-600 font-bold mt-1">￥99.00</p>
                    <router-link to="/products" class="btn-primary btn-sm mt-2 inline-block">查看更多</router-link>
                </div>
            </div>

            <div class="text-center mt-8">
                <router-link to="/products" class="btn-primary">浏览全部商品</router-link>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import Layout from '../components/Layout.vue'
</script>
