<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的优惠券</h2>

            <div class="flex items-center gap-2">
                <button @click="tab = 'available'" :class="tab === 'available' ? 'btn-primary' : 'btn-secondary'">可使用</button>
                <button @click="tab = 'used'" :class="tab === 'used' ? 'btn-primary' : 'btn-secondary'">已使用</button>
                <button @click="tab = 'expired'" :class="tab === 'expired' ? 'btn-primary' : 'btn-secondary'">已过期</button>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">加载中...</div>
            <div v-else-if="filteredCoupons.length === 0" class="card p-8 text-center text-gray-500">暂无优惠券</div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div v-for="c in filteredCoupons" :key="c.id" class="card overflow-hidden">
                    <div class="p-4 bg-gradient-to-r from-orange-50 to-yellow-50">
                        <div class="flex items-baseline gap-1 text-orange-600">
                            <span class="text-3xl font-bold">{{ formatDiscount(c) }}</span>
                        </div>
                        <div class="text-xs text-gray-600 mt-1">
                            满 {{ c.min_amount || 0 }} 元可用
                        </div>
                    </div>
                    <div class="p-3 text-sm">
                        <div class="font-semibold text-gray-800">{{ c.name || '优惠券' }}</div>
                        <div class="text-xs text-gray-500 mt-1">有效期至: {{ c.expire_time || '长期有效' }}</div>
                        <div class="text-xs text-gray-500">{{ c.description || '全场通用' }}</div>
                    </div>
                </div>
            </div>
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
