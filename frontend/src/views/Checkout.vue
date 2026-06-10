<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">结算</h2>

            <div v-if="loadingProduct && fromProduct" class="card p-8 text-center text-gray-500">
                加载商品信息...
            </div>
            <div v-else-if="!hasItems" class="card p-8 text-center text-gray-500">
                没有可以结算的商品。
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
                                    <option value="epay">在线支付</option>
                                    <option value="balance">余额支付</option>
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
                        <div class="card-body overflow-x-auto">
                            <table class="table w-full">
                                <thead>
                                    <tr>
                                        <th>商品</th>
                                        <th>单价</th>
                                        <th>数量</th>
                                        <th>小计</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="item in items" :key="item.product_id">
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
                            <span>{{ totalCount }} 件</span>
                        </div>
                        <div class="flex justify-between text-lg font-bold">
                            <span class="text-gray-800">合计</span>
                            <span class="text-blue-600">￥{{ totalPrice.toFixed(2) }}</span>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'
import { createOrder, payOrder } from '../api/order'
import { productDetail } from '../api/product'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const loading = ref(false)
const loadingProduct = ref(false)
const form = reactive({ email: '', pay_type: 'epay', remark: '' })
const productItem = ref(null)

const fromProduct = computed(() => !!route.query.product_id)

const items = computed(() => {
    if (productItem.value) {
        return [productItem.value]
    }
    return cartStore.items
})

const totalCount = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
const totalPrice = computed(() => items.value.reduce((sum, item) => sum + item.price * item.quantity, 0))
const hasItems = computed(() => items.value.length > 0)

function loadProduct() {
    const productId = route.query.product_id
    const quantity = parseInt(route.query.quantity) || 1
    if (!productId) return
    loadingProduct.value = true
    productDetail(productId).then((data) => {
        if (data && data.id) {
            productItem.value = {
                product_id: data.id,
                name: data.name,
                price: data.price,
                quantity: quantity
            }
        } else {
            alert('商品不存在')
            router.push('/products')
        }
    }).catch(() => {
        alert('商品不存在')
        router.push('/products')
    }).finally(() => {
        loadingProduct.value = false
    })
}

function submitOrder() {
    if (!form.email) {
        alert('请输入邮箱')
        return
    }
    const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRe.test(form.email)) {
        alert('邮箱格式不正确')
        return
    }
    if (items.value.length === 0) {
        alert('没有可结算的商品')
        return
    }
    loading.value = true
    const orderItem = items.value[0]
    createOrder({
        product_id: orderItem.product_id,
        quantity: orderItem.quantity,
        pay_type: form.pay_type,
        email: form.email,
        remark: form.remark
    }).then((data) => {
        if (!productItem.value) {
            cartStore.clearCart()
        }
        const orderNo = (data && data.order_no) ? data.order_no : null
        if (!orderNo) {
            router.push('/user/orders')
            return
        }
        return payOrder(orderNo).then((payData) => {
            if (payData && payData.pay_url) {
                window.location.href = payData.pay_url
            } else {
                router.push('/order/' + orderNo)
            }
        }).catch(() => {
            router.push('/order/' + orderNo)
        })
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    if (fromProduct.value) {
        loadProduct()
    }
})
</script>
