<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">购物车</h2>
            <div v-if="cartStore.items.length === 0" class="card p-8 text-center text-gray-500">
                购物车是空的，去选购商品吧。
                <div class="mt-4">
                    <router-link to="/products" class="btn-primary">去购物</router-link>
                </div>
            </div>
            <div v-else class="space-y-4">
                <div class="card">
                    <div class="card-header">
                        <div class="flex items-center justify-between">
                            <span class="font-semibold">商品列表</span>
                            <button @click="cartStore.clearCart()" class="btn-sm btn-danger">清空购物车</button>
                        </div>
                    </div>
                    <div class="card-body">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>商品</th>
                                    <th>单价</th>
                                    <th>数量</th>
                                    <th>小计</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="item in cartStore.items" :key="item.product_id">
                                    <td>{{ item.name }}</td>
                                    <td>￥{{ item.price }}</td>
                                    <td>
                                        <div class="flex items-center space-x-2">
                                            <button @click="cartStore.updateQuantity(item.product_id, item.quantity - 1)" class="btn-sm btn-secondary">-</button>
                                            <input v-model.number="item.quantity" type="number" min="1" class="form-input w-16 text-center" @change="cartStore.updateQuantity(item.product_id, item.quantity)">
                                            <button @click="cartStore.updateQuantity(item.product_id, item.quantity + 1)" class="btn-sm btn-secondary">+</button>
                                        </div>
                                    </td>
                                    <td>￥{{ (item.price * item.quantity).toFixed(2) }}</td>
                                    <td>
                                        <button @click="cartStore.removeItem(item.product_id)" class="btn-sm btn-danger">删除</button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body flex items-center justify-between">
                        <div class="text-lg">
                            共 <span class="font-semibold">{{ cartStore.totalCount }}</span> 件商品，
                            合计: <span class="text-blue-600 font-bold text-xl">￥{{ cartStore.totalPrice.toFixed(2) }}</span>
                        </div>
                        <router-link to="/checkout" class="btn-primary">去结算</router-link>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'

const cartStore = useCartStore()
</script>
