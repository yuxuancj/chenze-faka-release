<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的积分</h2>

            <el-card shadow="hover">
                <div class="flex items-center justify-between">
                    <div>
                        <div class="text-sm text-gray-500">当前积分</div>
                        <div class="text-4xl font-bold text-purple-700 mt-2">{{ currentPoints }}</div>
                    </div>
                    <el-button type="primary" size="large" @click="handleSignin" :disabled="signingIn || signedToday">
                        <span v-if="signingIn">签到中...</span>
                        <span v-else-if="signedToday">今日已签到</span>
                        <span v-else>每日签到</span>
                    </el-button>
                </div>
                <div v-if="signinReward" class="text-sm text-green-600 mt-3">
                    签到奖励: +{{ signinReward }} 积分
                </div>
            </el-card>

            <el-card shadow="hover">
                <template #header>
                    <span class="font-semibold">积分明细</span>
                </template>
                <el-table :data="logs" v-loading="loading" style="width: 100%">
                    <el-table-column prop="type_name" label="类型">
                        <template #default="scope">
                            {{ scope.row.type_name || (scope.row.type === 1 ? '获取' : (scope.row.type === 2 ? '消耗' : '其他')) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="变动">
                        <template #default="scope">
                            <span :class="(scope.row.change || scope.row.points || 0) >= 0 ? 'text-green-600' : 'text-red-600'" class="font-semibold">
                                {{ (scope.row.change || scope.row.points || 0) > 0 ? '+' : '' }}{{ scope.row.change || scope.row.points || 0 }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注">
                        <template #default="scope">
                            {{ scope.row.remark || scope.row.description || '-' }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="created_at" label="时间">
                        <template #default="scope">
                            {{ scope.row.created_at || '-' }}
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty description="暂无积分记录"></el-empty>
                    </template>
                </el-table>
            </el-card>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { pointsLogs, signIn } from '../api/points'
import { profile } from '../api/user'

const logs = ref([])
const currentPoints = ref(0)
const signedToday = ref(false)
const signinReward = ref(0)
const signingIn = ref(false)
const loading = ref(true)

function loadAll() {
    loading.value = true
    Promise.all([
        profile().then((d) => {
            currentPoints.value = (d && d.points) ? d.points : 0
        }).catch(() => {}),
        pointsLogs(1, 30).then((d) => {
            logs.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        }).catch(() => {})
    ]).finally(() => {
        loading.value = false
    })
}

function handleSignin() {
    if (signedToday.value) return
    signingIn.value = true
    signIn().then((d) => {
        signinReward.value = (d && d.points) ? d.points : ((d && d.reward) ? d.reward : 10)
        signedToday.value = true
        currentPoints.value += signinReward.value
        loadAll()
    }).catch(() => {}).finally(() => {
        signingIn.value = false
    })
}

onMounted(loadAll)
</script>
