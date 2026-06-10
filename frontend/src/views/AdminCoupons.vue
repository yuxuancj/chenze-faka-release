<template>
    <AdminLayout page-title="优惠券管理">
        <div class="mb-4 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-800">优惠券列表</h2>
            <el-button type="primary" size="small" @click="openCreateDialog">新增优惠券</el-button>
        </div>

        <el-table v-loading="loading" :data="coupons" style="width: 100%" border stripe empty-text="暂无数据">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" />
            <el-table-column label="类型" width="100">
                <template #default="scope">{{ scope.row.type === 2 ? '折扣券' : '满减券' }}</template>
            </el-table-column>
            <el-table-column label="面额/折扣" width="120">
                <template #default="scope">
                    {{ scope.row.type === 2 ? (scope.row.discount_rate || scope.row.discount_amount) + '折' : '￥' + (scope.row.discount_amount || 0) }}
                </template>
            </el-table-column>
            <el-table-column label="最低消费" width="120">
                <template #default="scope">￥{{ scope.row.min_amount || 0 }}</template>
            </el-table-column>
            <el-table-column prop="total_count" label="数量" width="100" />
            <el-table-column prop="used_count" label="已领取" width="100" />
            <el-table-column label="状态" width="100">
                <template #default="scope">
                    <el-tag v-if="scope.row.status === 1" type="success">启用</el-tag>
                    <el-tag v-else type="info">禁用</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="expire_time" label="有效期" width="180" />
            <el-table-column label="操作" width="180" fixed="right">
                <template #default="scope">
                    <el-button type="success" size="small" @click="openEditDialog(scope.row)">编辑</el-button>
                    <el-button type="danger" size="small" @click="deleteCoupon(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑优惠券' : '新增优惠券'" width="600px">
            <el-form :model="cform" label-width="120px">
                <el-form-item label="优惠券名称">
                    <el-input v-model="cform.name" />
                </el-form-item>
                <el-form-item label="类型">
                    <el-select v-model="cform.type" style="width: 100%">
                        <el-option label="满减券" :value="1" />
                        <el-option label="折扣券" :value="2" />
                    </el-select>
                </el-form-item>
                <el-form-item label="面额/折扣率">
                    <el-input-number v-model="cform.discount_amount" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="最低消费">
                    <el-input-number v-model="cform.min_amount" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="发放数量">
                    <el-input-number v-model="cform.total_count" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="有效期至">
                    <el-date-picker
                        v-model="cform.expire_time"
                        type="datetime"
                        placeholder="选择日期时间"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        style="width: 100%"
                    />
                </el-form-item>
                <el-form-item label="描述">
                    <el-input v-model="cform.description" type="textarea" :rows="2" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" :loading="saving" @click="submitCoupon">保存</el-button>
            </template>
        </el-dialog>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminCouponList, adminCouponCreate, adminCouponUpdate, adminCouponDelete } from '../api/coupon'

const coupons = ref([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(0)
const cform = reactive({
    name: '', type: 1, discount_amount: 0, min_amount: 0,
    total_count: 100, expire_time: '', description: ''
})

function resetForm() {
    cform.name = ''
    cform.type = 1
    cform.discount_amount = 0
    cform.min_amount = 0
    cform.total_count = 100
    cform.expire_time = ''
    cform.description = ''
}

function loadList() {
    loading.value = true
    adminCouponList(1, 100).then((d) => {
        coupons.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
    }).catch(() => {}).finally(() => { loading.value = false })
}

function openCreateDialog() {
    isEditing.value = false
    editingId.value = 0
    resetForm()
    dialogVisible.value = true
}

function openEditDialog(c) {
    isEditing.value = true
    editingId.value = c.id
    cform.name = c.name || ''
    cform.type = c.type || 1
    cform.discount_amount = c.type === 2 ? (c.discount_rate || 0) : (c.discount_amount || 0)
    cform.min_amount = c.min_amount || 0
    cform.total_count = c.total_count || 0
    cform.expire_time = c.expire_time || ''
    cform.description = c.description || ''
    dialogVisible.value = true
}

function closeDialog() {
    dialogVisible.value = false
    isEditing.value = false
    editingId.value = 0
    resetForm()
}

function submitCoupon() {
    if (!cform.name) {
        ElMessage.warning('请输入名称')
        return
    }
    saving.value = true
    const payload = {
        name: cform.name,
        type: cform.type,
        discount_amount: cform.type === 1 ? cform.discount_amount : 0,
        discount_rate: cform.type === 2 ? cform.discount_amount : 0,
        min_amount: cform.min_amount || 0,
        total_count: cform.total_count || 0,
        expire_time: cform.expire_time,
        description: cform.description,
        status: 1
    }
    const promise = isEditing.value && editingId.value
        ? adminCouponUpdate(editingId.value, payload)
        : adminCouponCreate(payload)
    promise.then(() => {
        ElMessage.success('保存成功')
        closeDialog()
        loadList()
    }).catch(() => {}).finally(() => { saving.value = false })
}

function deleteCoupon(id) {
    ElMessageBox.confirm('确定删除该优惠券？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        adminCouponDelete(id).then(() => {
            ElMessage.success('删除成功')
            loadList()
        }).catch(() => {})
    }).catch(() => {})
}

onMounted(loadList)
</script>
