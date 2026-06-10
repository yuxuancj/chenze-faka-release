<template>
    <AdminLayout page-title="分销设置">
        <el-card v-loading="loading" shadow="never">
            <template #header>
                <span class="font-semibold">佣金设置</span>
            </template>
            <el-form :model="settings" label-width="150px">
                <el-form-item label="一级佣金比例（%）">
                    <el-input-number v-model="settings.level1_rate" :precision="2" :step="0.01" :min="0" :max="100" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="二级佣金比例（%）">
                    <el-input-number v-model="settings.level2_rate" :precision="2" :step="0.01" :min="0" :max="100" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="三级佣金比例（%）">
                    <el-input-number v-model="settings.level3_rate" :precision="2" :step="0.01" :min="0" :max="100" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="最低提现金额">
                    <el-input-number v-model="settings.min_withdraw" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="功能开关">
                    <el-switch v-model="settings.enabled" :active-value="1" :inactive-value="0" active-text="启用分销" inactive-text="禁用分销" />
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
import { adminDistributionGet, adminDistributionSet } from '../api/distribution'

const settings = reactive({
    level1_rate: 10, level2_rate: 5, level3_rate: 2,
    min_withdraw: 10, enabled: 1
})
const saving = ref(false)
const loading = ref(true)

function load() {
    loading.value = true
    adminDistributionGet().then((d) => {
        if (d) {
            settings.level1_rate = d.level1_rate || d.level_1 || 10
            settings.level2_rate = d.level2_rate || d.level_2 || 5
            settings.level3_rate = d.level3_rate || d.level_3 || 2
            settings.min_withdraw = d.min_withdraw || d.min_amount || 10
            settings.enabled = d.enabled === undefined ? 1 : d.enabled
        }
    }).catch(() => {}).finally(() => { loading.value = false })
}

function save() {
    saving.value = true
    adminDistributionSet({ ...settings }).then(() => {
        ElMessage.success('保存成功')
    }).catch(() => {}).finally(() => { saving.value = false })
}

onMounted(load)
</script>
