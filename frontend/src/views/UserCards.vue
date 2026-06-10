<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的卡密</h2>

            <el-card v-if="loading" shadow="never" class="text-center">
                加载中...
            </el-card>
            <el-card v-else-if="cards.length === 0" shadow="never" class="text-center">
                <div>暂无已购买的卡密</div>
                <div class="mt-3">
                    <router-link to="/products">
                        <el-button type="primary" size="small">去购物</el-button>
                    </router-link>
                </div>
            </el-card>
            <el-row v-else :gutter="20">
                <el-col v-for="item in cards" :key="item.order_no + '-' + item.idx" :xs="24" :md="12" :lg="8" style="margin-bottom: 16px;">
                    <el-card shadow="hover">
                        <template #header>
                            <div class="flex items-center justify-between">
                                <span class="font-semibold truncate" style="max-width: 240px;">{{ item.product_name || item.product || '商品' }}</span>
                                <span class="text-xs text-gray-500">{{ item.order_no }}</span>
                            </div>
                        </template>
                        <div class="text-xs text-gray-500 mb-2">购买时间：{{ item.created_at || '-' }}</div>
                        <div v-for="(card, idx2) in item.cards" :key="idx2" class="flex items-center gap-2 mb-2">
                            <el-input
                                :model-value="typeof card === 'string' ? card : (card.card_data || card.card_no || '')"
                                readonly
                                size="small"
                                class="font-mono text-xs flex-1"
                            />
                            <el-button type="primary" size="small" @click="copyCard(typeof card === 'string' ? card : (card.card_data || card.card_no || ''))">
                                复制
                            </el-button>
                        </div>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Layout from '../components/Layout.vue'
import { orderList, orderDetail } from '../api/order'

const cards = ref([])
const loading = ref(true)

function copyCard(data) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(data).then(() => {
            ElMessage.success('已复制卡密')
        }).catch(() => {
            fallbackCopy(data)
        })
    } else {
        fallbackCopy(data)
    }
}

function fallbackCopy(data) {
    const ta = document.createElement('textarea')
    ta.value = data
    document.body.appendChild(ta)
    ta.select()
    try {
        document.execCommand('copy')
        ElMessage.success('已复制卡密')
    } catch (e) {
        ElMessage.error('复制失败，请手动复制')
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
