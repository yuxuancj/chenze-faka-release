<template>
    <AdminLayout page-title="秒杀活动管理">
        <div class="mb-4 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-800">秒杀活动列表</h2>
            <el-button type="primary" size="small" @click="openCreateDialog">新增秒杀</el-button>
        </div>

        <el-table v-loading="loading" :data="seckills" style="width: 100%" border stripe empty-text="暂无数据">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_id" label="商品ID" width="100" />
            <el-table-column prop="sku_id" label="SKU ID" width="100" />
            <el-table-column label="秒杀价" width="120">
                <template #default="scope" style="color: #dc2626; font-weight: 600">￥{{ scope.row.seckill_price }}</template>
            </el-table-column>
            <el-table-column prop="stock" label="库存" width="100" />
            <el-table-column prop="sold" label="已售" width="100" />
            <el-table-column prop="limit_per_user" label="限购" width="100" />
            <el-table-column label="状态" width="100">
                <template #default="scope">
                    <el-tag v-if="scope.row.status === 1" type="success">启用</el-tag>
                    <el-tag v-else type="info">禁用</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="时间" width="320">
                <template #default="scope">{{ formatDateTime(scope.row.start_time) }} - {{ formatDateTime(scope.row.end_time) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
                <template #default="scope">
                    <el-button type="success" size="small" @click="openEditDialog(scope.row)">编辑</el-button>
                    <el-button type="danger" size="small" @click="deleteSeckill(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑秒杀' : '新增秒杀'" width="600px">
            <el-form :model="sform" label-width="120px">
                <el-form-item label="商品ID">
                    <el-input-number v-model="sform.product_id" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="SKU ID">
                    <el-input-number v-model="sform.sku_id" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="秒杀价">
                    <el-input-number v-model="sform.seckill_price" :precision="2" :step="0.01" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="秒杀库存">
                    <el-input-number v-model="sform.stock" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="限购数量">
                    <el-input-number v-model="sform.limit_per_user" :step="1" :min="0" controls-position="right" style="width: 100%" />
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="sform.status" style="width: 100%">
                        <el-option label="启用" :value="1" />
                        <el-option label="禁用" :value="0" />
                    </el-select>
                </el-form-item>
                <el-form-item label="开始时间">
                    <el-date-picker
                        v-model="sform.start_time"
                        type="datetime"
                        placeholder="选择日期时间"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        style="width: 100%"
                    />
                </el-form-item>
                <el-form-item label="结束时间">
                    <el-date-picker
                        v-model="sform.end_time"
                        type="datetime"
                        placeholder="选择日期时间"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        style="width: 100%"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" :loading="saving" @click="submitSeckill">保存</el-button>
            </template>
        </el-dialog>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminSeckillList, adminSeckillCreate, adminSeckillUpdate, adminSeckillDelete } from '../api/seckill'

function formatDateTime(dateStr) {
    if (!dateStr) return '-'
    try {
        const d = new Date(dateStr)
        if (isNaN(d.getTime())) return '-'
        return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    } catch (e) {
        return '-'
    }
}

const seckills = ref([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(0)
const sform = reactive({
    product_id: 0, sku_id: 0, seckill_price: 0, stock: 0, limit_per_user: 1,
    status: 1, start_time: '', end_time: ''
})

function resetForm() {
    sform.product_id = 0
    sform.sku_id = 0
    sform.seckill_price = 0
    sform.stock = 0
    sform.limit_per_user = 1
    sform.status = 1
    sform.start_time = ''
    sform.end_time = ''
}

function loadList() {
    loading.value = true
    adminSeckillList(1, 100).then((d) => {
        seckills.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
    }).catch(() => {}).finally(() => { loading.value = false })
}

function openCreateDialog() {
    isEditing.value = false
    editingId.value = 0
    resetForm()
    dialogVisible.value = true
}

function openEditDialog(s) {
    isEditing.value = true
    editingId.value = s.id
    sform.product_id = s.product_id || 0
    sform.sku_id = s.sku_id || 0
    sform.seckill_price = s.seckill_price || 0
    sform.stock = s.stock || 0
    sform.limit_per_user = s.limit_per_user || 1
    sform.status = s.status || 1
    sform.start_time = s.start_time || ''
    sform.end_time = s.end_time || ''
    dialogVisible.value = true
}

function closeDialog() {
    dialogVisible.value = false
    isEditing.value = false
    editingId.value = 0
    resetForm()
}

function submitSeckill() {
    if (!sform.product_id) {
        ElMessage.warning('请输入商品ID')
        return
    }
    saving.value = true
    const payload = { ...sform }
    const promise = isEditing.value && editingId.value
        ? adminSeckillUpdate(editingId.value, payload)
        : adminSeckillCreate(payload)
    promise.then(() => {
        ElMessage.success('保存成功')
        closeDialog()
        loadList()
    }).catch(() => {}).finally(() => { saving.value = false })
}

function deleteSeckill(id) {
    ElMessageBox.confirm('确定删除？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        adminSeckillDelete(id).then(() => {
            ElMessage.success('删除成功')
            loadList()
        }).catch(() => {})
    }).catch(() => {})
}

onMounted(loadList)
</script>
