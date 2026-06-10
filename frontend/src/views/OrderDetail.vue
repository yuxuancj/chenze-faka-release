<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">订单详情</h2>
            <div v-if="order.order_no" class="space-y-4">
                <div class="card">
                    <div class="card-header font-semibold">订单信息</div>
                    <div class="card-body">
                        <table class="table">
                            <tbody>
                                <tr><td class="w-32 text-gray-500">订单号</td><td>{{ order.order_no }}</td></tr>
                                <tr><td class="text-gray-500">商品名称</td><td>{{ order.product_name }}</td></tr>
                                <tr><td class="text-gray-500">购买数量</td><td>{{ order.quantity }}</td></tr>
                                <tr><td class="text-gray-500">订单金额</td><td>￥{{ order.amount }}</td></tr>
                                <tr><td class="text-gray-500">状态</td>
                                    <td>
                                        <span v-if="order.status === 1" class="badge-yellow">待支付</span>
                                        <span v-else-if="order.status === 2" class="badge-green">已完成</span>
                                        <span v-else-if="order.status === 3" class="badge-red">已取消</span>
                                        <span v-else class="badge-gray">未知</span>
                                    </td>
                                </tr>
                                <tr><td class="text-gray-500">邮箱</td><td>{{ order.email }}</td></tr>
                                <tr><td class="text-gray-500">创建时间</td><td>{{ order.created_at }}</td></tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div v-if="order.status === 1" class="card">
                    <div class="card-body flex items-center justify-between">
                        <span>待支付金额: <span class="text-blue-600 font-bold text-xl">￥{{ order.amount }}</span></span>
                        <button @click="payNow" :disabled="paying" class="btn-primary">
                            {{ paying ? '支付中...' : '立即支付' }}
                        </button>
                    </div>
                </div>
                <div v-if="cards && cards.length > 0" class="card">
                    <div class="card-header font-semibold">卡密信息</div>
                    <div class="card-body">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>卡号</th>
                                    <th>密码</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(card, idx) in cards" :key="idx">
                                    <td>{{ card.card_no }}</td>
                                    <td>{{ card.card_pwd }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <div v-else class="card p-8 text-center text-gray-500">
                加载中...
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Layout from '../components/Layout.vue'
import { orderDetail, payOrder } from '../api/order'

const route = useRoute()
const order = ref({})
const cards = ref([])
const paying = ref(false)

function loadDetail() {
    orderDetail(route.params.order_no).then((res) => {
        order.value = res.data || {}
        cards.value = (res.data && res.data.cards) ? res.data.cards : []
    }).catch(() => {})
}

function payNow() {
    paying.value = true
    payOrder(order.value.order_no).then(() => {
        alert('支付成功')
        loadDetail()
    }).catch(() => {}).finally(() => {
        paying.value = false
    })
}

onMounted(loadDetail)
</script>
