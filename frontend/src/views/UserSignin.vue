<template>
    <Layout>
        <div class="max-w-md mx-auto space-y-4">
            <el-card class="gradient-header" body-style="padding: 0">
                <div class="px-6 py-8 text-center">
                    <h1 class="text-2xl font-bold text-white mb-2">每日签到</h1>
                    <p class="text-green-100">连续签到领取额外积分奖励</p>
                </div>
            </el-card>

            <el-card v-loading="loading" shadow="hover">
                <div class="text-center space-y-4">
                    <div class="text-5xl">{{ canSignin ? '🎯' : '✅' }}</div>
                    <div v-if="canSignin" class="space-y-4">
                        <p class="text-gray-600">
                            今日签到可获得 <span class="text-green-600 font-bold text-xl">{{ signinReward }}</span> 积分
                        </p>
                        <p v-if="continuousDays > 0" class="text-sm text-orange-500">
                            继续坚持！已连续签到 <strong>{{ continuousDays }}</strong> 天
                        </p>
                        <p class="text-xs text-gray-400">
                            连续3天额外+{{ continuousReward }}积分，连续7天双倍奖励
                        </p>
                        <el-button type="primary" size="large" class="w-full" @click="doSignin" :loading="signing">
                            立即签到
                        </el-button>
                    </div>
                    <div v-else class="space-y-3">
                        <p class="text-gray-600 text-lg">今日已签到</p>
                        <el-tag type="success" size="large">🔥 连续签到 {{ continuousDays }} 天</el-tag>
                        <p class="text-green-600 font-medium">明日再来领取更多积分</p>
                    </div>
                </div>
            </el-card>

            <el-card shadow="hover">
                <template #header>
                    <span class="font-semibold">签到规则</span>
                </template>
                <div class="space-y-3 text-sm text-gray-600">
                    <div class="flex items-start gap-2">
                        <el-icon class="text-green-500"><Check /></el-icon>
                        <span>每日签到可获得 <strong>{{ signinReward }}</strong> 积分</span>
                    </div>
                    <div class="flex items-start gap-2">
                        <el-icon class="text-orange-500"><Promotion /></el-icon>
                        <span>连续签到{{ continuousRewardDays }}天额外奖励 <strong>+{{ continuousReward }}</strong> 积分</span>
                    </div>
                    <div class="flex items-start gap-2">
                        <el-icon class="text-purple-500"><Star /></el-icon>
                        <span>连续签到{{ doubleRewardDays }}天奖励<strong>翻倍</strong></span>
                    </div>
                    <div class="flex items-start gap-2">
                        <el-icon class="text-blue-500"><Money /></el-icon>
                        <span>积分可抵扣订单金额</span>
                    </div>
                </div>
            </el-card>

            <el-card shadow="hover">
                <template #header>
                    <div class="flex justify-between items-center">
                        <span class="font-semibold">我的积分</span>
                        <router-link to="/user/points" class="text-blue-500 text-sm font-normal">查看全部</router-link>
                    </div>
                </template>
                <div class="text-center py-4">
                    <span class="text-5xl font-bold text-green-600">{{ balance || 0 }}</span>
                    <p class="text-sm text-gray-500 mt-1">可用积分</p>
                </div>
                <el-divider />
                <p class="text-sm text-gray-500 mb-2">最近积分变动</p>
                <el-table :data="logs.slice(0, 5)" v-loading="logsLoading" size="small" :show-header="false" style="width: 100%">
                    <el-table-column>
                        <template #default="scope">
                            <div class="flex justify-between items-center text-sm py-2">
                                <div>
                                    <span class="text-gray-700">{{ scope.row.description || getTypeDesc(scope.row.type) }}</span>
                                    <span v-if="scope.row.created_at" class="text-xs text-gray-400 block">
                                        {{ formatTime(scope.row.created_at) }}
                                    </span>
                                </div>
                                <span :class="(scope.row.amount || scope.row.change || 0) > 0 ? 'text-green-600 font-medium' : 'text-red-600'">
                                    {{ (scope.row.amount || scope.row.change || 0) > 0 ? '+' : '' }}{{ scope.row.amount || scope.row.change || 0 }}
                                </span>
                            </div>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty description="暂无记录" :image-size="60"></el-empty>
                    </template>
                </el-table>
            </el-card>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Check, Promotion, Star, Money } from '@element-plus/icons-vue'
import Layout from '../components/Layout.vue'
import { signIn } from '../api/points'
import { ElMessage } from 'element-plus'

const loading = ref(true)
const logsLoading = ref(false)
const signing = ref(false)
const canSignin = ref(false)
const continuousDays = ref(0)
const balance = ref(0)
const logs = ref([])
const signinReward = ref(10)
const continuousReward = ref(5)
const continuousRewardDays = ref(3)
const doubleRewardDays = ref(7)

function getTypeDesc(type) {
    const map = {
        'signin': '每日签到',
        'order': '购物返积分',
        'deduct': '积分抵扣',
        'refund': '退款返还'
    }
    return map[type] || '积分变动'
}

function formatTime(timeStr) {
    if (!timeStr) return ''
    const d = new Date(timeStr)
    const now = new Date()
    const diff = now - d
    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
    if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
    return d.toLocaleDateString()
}

async function loadStatus() {
    loading.value = true
    try {
        const resp = await fetch('/api/v1/user/signin/status')
        const d = await resp.json()
        if (d.code === 0) {
            canSignin.value = !!d.data?.can_signin
            continuousDays.value = d.data?.continuous_days || 0
            balance.value = d.data?.balance || 0
        } else {
            canSignin.value = true
        }
    } catch (e) {
        console.error('加载签到状态失败:', e)
        canSignin.value = true
    } finally {
        loading.value = false
    }
}

async function loadLogs() {
    logsLoading.value = true
    try {
        const resp = await fetch('/api/v1/user/points/logs?page=1&size=10')
        const d = await resp.json()
        if (d.code === 0) {
            logs.value = (d.data?.list) ? d.data.list : []
        }
    } catch (e) {
        console.error('加载积分记录失败:', e)
    } finally {
        logsLoading.value = false
    }
}

async function doSignin() {
    signing.value = true
    try {
        const result = await signIn()
        const data = result || {}
        const earned = data.points || signinReward.value
        const days = data.continuous_days || (continuousDays.value + 1)

        ElMessage.success(`签到成功！获得 ${earned} 积分`)
        canSignin.value = false
        continuousDays.value = days
        balance.value += earned

        await loadLogs()
    } catch (e) {
        const msg = e.message || ''
        if (msg.includes('已签到')) {
            ElMessage.info('今日已签到')
            canSignin.value = false
        } else if (msg.includes('登录')) {
            ElMessage.error('请先登录')
            window.location.href = '/user/login'
        } else {
            ElMessage.error(msg || '签到失败，请重试')
        }
    } finally {
        signing.value = false
    }
}

onMounted(() => {
    loadStatus()
    loadLogs()
})
</script>

<style scoped>
.gradient-header :deep(.el-card__body) {
    padding: 0;
}
.gradient-header {
    background: linear-gradient(to right, #10b981, #f59e0b);
    border: none;
}
</style>
