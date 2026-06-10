<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">结算</h2>
            <div v-if="cartStore.items.length === 0" class="card p-8 text-center text-gray-500">
                购物车是空的。
                <div class="mt-4">
                    <router-link to="/products" class="btn-primary">去购物</router-link>
                </div>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div class="md:col-span-2 space-y-4">
                    <div class="card">
                        <div class="card-header font-semibold">订单信息</div>
                        <div class="card-body space-y-4">
                            <div>
                                <label class="form-label">联系邮箱</label>
                                <input v-model="form.email" type="email" class="form-input" placeholder="请输入接收卡密的邮箱">
                            </div>
                            <div>
                                <label class="form-label">支付方式</label>
                                <select v-model="form.pay_type" class="form-input">
                                    <option value="1">在线支付</option>
                                    <option value="2">余额支付</option>
                                </select>
                            </div>
                            <div>
                                <label class="form-label">备注</label>
                                <textarea v-model="form.remark" class="form-input" rows="3" placeholder="选填"></textarea>
                            </div>
                        </div>
                    </div>
                    <div class="card">
                        <div class="card-header font-semibold">商品清单</div>
                        <div class="card-body">
                            <table class="table">
                                <thead>
                                    <tr>
                                        <th>商品</th>
                                        <th>单价</th>
                                        <th>数量</th>
                                        <th>小计</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="item in cartStore.items" :key="item.product_id">
                                        <td>{{ item.name }}</td>
                                        <td>￥{{ item.price }}</td>
                                        <td>{{ item.quantity }}</td>
                                        <td>￥{{ (item.price * item.quantity).toFixed(2) }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div class="card h-fit">
                    <div class="card-header font-semibold">结算</div>
                    <div class="card-body space-y-3">
                        <div class="flex justify-between text-sm">
                            <span class="text-gray-600">商品总数</span>
                            <span>{{ cartStore.totalCount }} 件</span>
                        </div>
                        <div class="flex justify-between text-lg font-bold">
                            <span class="text-gray-800">合计</span>
                            <span class="text-blue-600">￥{{ cartStore.totalPrice.toFixed(2) }}</span>
                        </div>
                        <button @click="submitOrder" :disabled="loading" class="btn-primary w-full">
                            {{ loading ? '提交中...' : '提交订单' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'
import { createOrder } from '../api/order'

const router = useRouter()
const cartStore = useCartStore()
const loading = ref(false)
const form = reactive({ email: '', pay_type: '1', remark: '' })

function submitOrder() {
    if (!form.email) {
        alert('请输入邮箱')
        return
    }
    loading.value = true
    const items = cartStore.items
    const item = items[0]
    createOrder({
        product_id: item.product_id,
        quantity: item.quantity,
        pay_type: form.pay_type,
        email: form.email,
        remark: form.remark
    }).then((res) => {
        cartStore.clearCart()
        if (res.data && res.data.order_no) {
            router.push('/order/' + res.data.order_no)
        } else {
            router.push('/user/orders')
        }
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}
</script>
