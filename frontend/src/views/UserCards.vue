<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的卡密</h2>

            <div v-if="loading" class="card p-8 text-center text-gray-500">加载中...</div>
            <div v-else-if="cards.length === 0" class="card p-8 text-center text-gray-500">
                暂无已购买的卡密
                <div class="mt-3">
                    <router-link to="/products" class="btn-primary btn-sm">去购物</router-link>
                </div>
            </div>
            <div v-else class="space-y-3">
                <div v-for="item in cards" :key="item.order_no + '-' + item.idx" class="card">
                    <div class="card-header flex items-center justify-between">
                        <span class="font-semibold">{{ item.product_name || item.product || '商品' }}</span>
                        <span class="text-xs text-gray-500">订单: {{ item.order_no }}</span>
                    </div>
                    <div class="card-body space-y-2">
                        <div class="text-xs text-gray-500">购买时间: {{ item.created_at || '-' }}</div>
                        <div v-for="(card, idx2) in item.cards" :key="idx2" class="flex items-center gap-2">
                            <input :value="typeof card === 'string' ? card : (card.card_data || card.card_no || '')" readonly class="form-input flex-1 font-mono text-xs">
                            <button @click="copyCard(typeof card === 'string' ? card : (card.card_data || card.card_no || ''))" class="btn-sm btn-secondary">复制</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { orderList, orderDetail } from '../api/order'

const cards = ref([])
const loading = ref(true)

function copyCard(data) {
    const ta = document.createElement('textarea')
    ta.value = data
    document.body.appendChild(ta)
    ta.select()
    try {
        document.execCommand('copy')
        alert('已复制卡密')
    } catch (e) {
        alert('复制失败，请手动复制')
    }
    document.body.removeChild(ta)
}

function load() {
    loading.value = true
    orderList(1, 50).then((d) => {
        const list = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        const paidOrders = list.filter((o) => o.status === 1 || o.status === 2)
        const tasks = paidOrders.map((o) => {
            return orderDetail(o.order_no).then((detail) => {
                let cs = []
                if (detail && detail.cards) cs = detail.cards
                else if (detail && detail.order && detail.order.cards) cs = detail.order.cards
                if (cs && cs.length > 0) {
                    cards.value.push({
                        order_no: o.order_no,
                        product_name: o.product_snapshot || o.product_name || '',
                        created_at: o.created_at,
                        idx: cards.value.length,
                        cards: cs
                    })
                }
            }).catch(() => {})
        })
        Promise.all(tasks).finally(() => {
            loading.value = false
        })
    }).catch(() => {
        loading.value = false
    })
}

onMounted(load)
</script>
