<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的优惠券</h2>

            <el-tabs v-model="tab">
                <el-tab-pane label="可使用" name="available"></el-tab-pane>
                <el-tab-pane label="已使用" name="used"></el-tab-pane>
                <el-tab-pane label="已过期" name="expired"></el-tab-pane>
            </el-tabs>

            <el-empty v-if="loading" description="加载中..."></el-empty>
            <el-empty v-else-if="filteredCoupons.length === 0" description="暂无优惠券"></el-empty>
            <el-row v-else :gutter="16">
                <el-col v-for="c in filteredCoupons" :key="c.id" :xs="24" :sm="12" :md="8">
                    <el-card class="coupon-card mb-4" shadow="hover">
                        <template #header>
                            <div class="flex items-baseline gap-1 text-orange-600">
                                <span class="text-3xl font-bold">{{ formatDiscount(c) }}</span>
                            </div>
                            <div class="text-xs text-gray-600 mt-1">
                                满 {{ c.min_amount || 0 }} 元可用
                            </div>
                        </template>
                        <div class="text-sm">
                            <div class="font-semibold text-gray-800">{{ c.name || '优惠券' }}</div>
                            <div class="text-xs text-gray-500 mt-1">有效期至: {{ c.expire_time || '长期有效' }}</div>
                            <div class="text-xs text-gray-500">{{ c.description || '全场通用' }}</div>
                        </div>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { userCoupons } from '../api/coupon'

const coupons = ref([])
const loading = ref(true)
const tab = ref('available')

const filteredCoupons = computed(() => {
    return coupons.value.filter((c) => {
        if (tab.value === 'available') return c.status === 0 || c.status === undefined
        if (tab.value === 'used') return c.status === 1
        if (tab.value === 'expired') return c.status === 2
        return true
    })
})

function formatDiscount(c) {
    if (c.type === 2 || c.discount_rate) {
        return (c.discount_rate || c.value) + '折'
    }
    return '￥' + (c.discount_amount || c.value || 0)
}

function load() {
    loading.value = true
    userCoupons().then((d) => {
        coupons.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
    }).catch(() => { coupons.value = [] }).finally(() => {
        loading.value = false
    })
}

onMounted(load)
</script>

<style scoped>
.coupon-card :deep(.el-card__header) {
    background: linear-gradient(to right, #fff7ed, #fefce8);
    padding: 16px;
}
</style>
