<template>
    <AdminLayout page-title="积分设置">
        <el-card v-loading="loading" shadow="never">
            <template #header>
                <span class="font-semibold">积分规则</span>
            </template>
            <el-form :model="settings" label-width="180px">
                <el-form-item label="消费1元获得积分">
                    <el-input-number v-model="settings.points_per_yuan" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="积分抵扣比例(100积分=多少元)">
                    <el-input-number v-model="settings.yuan_per_hundred" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="每日签到奖励积分">
                    <el-input-number v-model="settings.signin_reward" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="连续签到奖励积分">
                    <el-input-number v-model="settings.continuous_reward" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="最高抵扣比例(%)">
                    <el-input-number v-model="settings.max_discount_rate" :step="1" :min="0" :max="100" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="功能开关">
                    <el-switch v-model="settings.enabled" :active-value="1" :inactive-value="0" active-text="启用积分" inactive-text="禁用积分" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" :loading="saving" @click="save">保存设置</el-button>
                </el-form-item>
            </el-form>
        </el-card>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminPointsGet, adminPointsSet } from '../api/points'

const settings = reactive({
    points_per_yuan: 1, yuan_per_hundred: 1, signin_reward: 10,
    continuous_reward: 5, max_discount_rate: 50, enabled: 1
})
const saving = ref(false)
const loading = ref(true)

function load() {
    loading.value = true
    adminPointsGet().then((d) => {
        if (d) {
            settings.points_per_yuan = d.points_per_yuan || d.points_per_rmb || 1
            settings.yuan_per_hundred = d.yuan_per_hundred || d.deduct_rate || 1
            settings.signin_reward = d.signin_reward || d.sign_in_points || 10
            settings.continuous_reward = d.continuous_reward || 5
            settings.max_discount_rate = d.max_discount_rate || d.max_rate || 50
            settings.enabled = d.enabled === undefined ? 1 : d.enabled
        }
    }).catch(() => {}).finally(() => { loading.value = false })
}

function save() {
    saving.value = true
    adminPointsSet({ ...settings }).then(() => {
        ElMessage.success('保存成功')
    }).catch(() => {}).finally(() => { saving.value = false })
}

onMounted(load)
</script>
