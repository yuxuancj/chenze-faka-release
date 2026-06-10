<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">分销中心</h2>

            <el-row v-if="!loading" :gutter="16">
                <el-col :xs="24" :sm="8">
                    <el-card shadow="hover">
                        <div class="text-sm text-gray-500">累计佣金</div>
                        <div class="text-2xl font-bold text-gray-800 mt-2">￥{{ summary.total_commission || 0 }}</div>
                    </el-card>
                </el-col>
                <el-col :xs="24" :sm="8">
                    <el-card shadow="hover">
                        <div class="text-sm text-gray-500">可提现金额</div>
                        <div class="text-2xl font-bold text-blue-600 mt-2">￥{{ summary.available_commission || 0 }}</div>
                    </el-card>
                </el-col>
                <el-col :xs="24" :sm="8">
                    <el-card shadow="hover">
                        <div class="text-sm text-gray-500">下级人数</div>
                        <div class="text-2xl font-bold text-gray-800 mt-2">{{ summary.team_count || 0 }}</div>
                    </el-card>
                </el-col>
            </el-row>

            <el-card v-if="!loading" shadow="hover">
                <template #header>
                    <span class="font-semibold">推广链接</span>
                </template>
                <div class="space-y-3">
                    <div class="flex items-center gap-2">
                        <el-input v-model="promoLink" readonly class="font-mono text-xs" :suffix-icon="CopyDocument" />
                        <el-button type="primary" @click="copyLink">复制</el-button>
                    </div>
                    <div v-if="posterUrl" class="mt-3">
                        <div class="text-sm text-gray-600 mb-2">推广海报</div>
                        <el-image :src="posterUrl" fit="cover" :preview-src-list="[posterUrl]" style="max-width: 300px" />
                    </div>
                    <el-button v-if="!posterLoading && !posterUrl" @click="loadPoster">
                        生成海报
                    </el-button>
                    <div v-if="posterLoading" class="text-sm text-gray-500">生成中...</div>
                </div>
            </el-card>

            <el-card v-if="!loading" shadow="hover">
                <template #header>
                    <span class="font-semibold">佣金明细</span>
                </template>
                <el-table :data="commissions" style="width: 100%">
                    <el-table-column prop="from_user_nickname" label="来源">
                        <template #default="scope">
                            {{ scope.row.from_user_nickname || ('用户' + (scope.row.from_user_id || '')) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="amount" label="金额">
                        <template #default="scope">
                            <span class="text-blue-600 font-semibold">￥{{ scope.row.amount }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="level" label="等级"></el-table-column>
                    <el-table-column label="状态">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">已结算</el-tag>
                            <el-tag v-else type="warning">待结算</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="created_at" label="时间">
                        <template #default="scope">
                            {{ scope.row.created_at || scope.row.settled_at || '-' }}
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty description="暂无佣金记录"></el-empty>
                    </template>
                </el-table>
            </el-card>

            <el-card v-if="!loading" shadow="hover">
                <template #header>
                    <span class="font-semibold">申请提现</span>
                </template>
                <el-form :inline="true" :model="withdrawForm">
                    <el-form-item label="提现金额">
                        <el-input-number v-model="withdrawForm.amount" :min="0" :step="0.01" />
                    </el-form-item>
                    <el-form-item label="收款方式">
                        <el-select v-model="withdrawForm.account_type" placeholder="请选择">
                            <el-option label="支付宝" value="alipay" />
                            <el-option label="微信" value="wechat" />
                            <el-option label="银行卡" value="bank" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="收款账号">
                        <el-input v-model="withdrawForm.account" placeholder="请输入收款账号" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="applyWithdraw" :loading="withdrawing">
                            提交申请
                        </el-button>
                    </el-form-item>
                </el-form>
            </el-card>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { CopyDocument } from '@element-plus/icons-vue'
import Layout from '../components/Layout.vue'
import { distributionSummary, distributionCommissions, distributionPoster, withdrawApply } from '../api/distribution'
import { ElMessage, ElMessageBox } from 'element-plus'

const summary = ref({})
const commissions = ref([])
const loading = ref(true)
const withdrawing = ref(false)
const posterUrl = ref('')
const posterLoading = ref(false)
const withdrawForm = reactive({ amount: 0, account_type: 'alipay', account: '' })

const promoLink = computed(() => {
    const code = summary.value.invite_code || ''
    const base = window.location.origin + '/register?invite='
    return base + code
})

function loadAll() {
    loading.value = true
    Promise.all([
        distributionSummary().then((d) => { summary.value = d || {} }).catch(() => {}),
        distributionCommissions(1, 20).then((d) => {
            commissions.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        }).catch(() => {})
    ]).finally(() => {
        loading.value = false
    })
}

function loadPoster() {
    posterLoading.value = true
    distributionPoster().then((d) => {
        posterUrl.value = (d && (d.poster_url || d.url || d.image)) || ''
    }).catch(() => {}).finally(() => {
        posterLoading.value = false
    })
}

function copyLink() {
    const ta = document.createElement('textarea')
    ta.value = promoLink.value
    document.body.appendChild(ta)
    ta.select()
    try {
        document.execCommand('copy')
        ElMessage.success('已复制推广链接')
    } catch (e) {
        ElMessage.error('复制失败，请手动复制')
    }
    document.body.removeChild(ta)
}

function applyWithdraw() {
    if (!withdrawForm.amount || withdrawForm.amount <= 0) {
        ElMessage.warning('请输入正确的提现金额')
        return
    }
    if (!withdrawForm.account) {
        ElMessage.warning('请输入收款账号')
        return
    }
    withdrawing.value = true
    withdrawApply(withdrawForm.amount, withdrawForm.account_type, withdrawForm.account).then(() => {
        ElMessage.success('申请已提交，等待审核')
        withdrawForm.amount = 0
        withdrawForm.account = ''
    }).catch(() => {}).finally(() => {
        withdrawing.value = false
    })
}

onMounted(loadAll)
</script>
