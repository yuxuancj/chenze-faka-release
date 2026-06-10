<template>
    <AdminLayout page-title="秒杀活动管理">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <h2 class="text-lg font-semibold text-gray-800">秒杀活动列表</h2>
                <button @click="showCreateForm = true" class="btn-primary btn-sm">新增秒杀</button>
            </div>

            <div v-if="showCreateForm || showEditForm" class="card">
                <div class="card-header font-semibold">{{ showEditForm ? '编辑秒杀' : '新增秒杀' }}</div>
                <div class="card-body space-y-3">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div>
                            <label class="form-label">商品ID</label>
                            <input v-model.number="sform.product_id" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">SKU ID (可选)</label>
                            <input v-model.number="sform.sku_id" type="number" class="form-input" placeholder="0 表示无 SKU">
                        </div>
                        <div>
                            <label class="form-label">秒杀价</label>
                            <input v-model.number="sform.price" type="number" step="0.01" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">秒杀库存</label>
                            <input v-model.number="sform.stock" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">限购数量</label>
                            <input v-model.number="sform.limit_per_user" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">状态</label>
                            <select v-model.number="sform.status" class="form-input">
                                <option :value="1">启用</option>
                                <option :value="0">禁用</option>
                            </select>
                        </div>
                        <div>
                            <label class="form-label">开始时间</label>
                            <input v-model="sform.start_time" type="datetime-local" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">结束时间</label>
                            <input v-model="sform.end_time" type="datetime-local" class="form-input">
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <button @click="submitSeckill" :disabled="saving" class="btn-primary btn-sm">
                            {{ saving ? '保存中...' : '保存' }}
                        </button>
                        <button @click="cancelForm" class="btn-secondary btn-sm">取消</button>
                    </div>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">加载中...</div>
            <div v-else class="card">
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>商品ID</th>
                                <th>SKU ID</th>
                                <th>秒杀价</th>
                                <th>库存</th>
                                <th>已售</th>
                                <th>限购</th>
                                <th>状态</th>
                                <th>时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="seckills.length === 0">
                                <td colspan="10" class="text-center text-gray-500 py-6">暂无数据</td>
                            </tr>
                            <tr v-for="s in seckills" :key="s.id">
                                <td>{{ s.id }}</td>
                                <td>{{ s.product_id }}</td>
                                <td>{{ s.sku_id || '-' }}</td>
                                <td class="text-red-600 font-semibold">￥{{ s.price }}</td>
                                <td>{{ s.stock }}</td>
                                <td>{{ s.sold || 0 }}</td>
                                <td>{{ s.limit_per_user }}</td>
                                <td>{{ s.status === 1 ? '启用' : '禁用' }}</td>
                                <td class="text-xs">{{ s.start_time }} - {{ s.end_time }}</td>
                                <td>
                                    <button @click="editSeckill(s)" class="btn-sm btn-secondary">编辑</button>
                                    <button @click="deleteSeckill(s.id)" class="btn-sm btn-danger ml-1">删除</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminSeckillList, adminSeckillCreate, adminSeckillUpdate, adminSeckillDelete } from '../api/seckill'

const seckills = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreateForm = ref(false)
const showEditForm = ref(false)
const editingId = ref(0)
const sform = reactive({
    product_id: 0, sku_id: 0, price: 0, stock: 0, limit_per_user: 1,
    status: 1, start_time: '', end_time: ''
})

function resetForm() {
    sform.product_id = 0; sform.sku_id = 0; sform.price = 0; sform.stock = 0
    sform.limit_per_user = 1; sform.status = 1; sform.start_time = ''; sform.end_time = ''
}

function loadList() {
    loading.value = true
    adminSeckillList(1, 100).then((d) => {
        seckills.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
    }).catch(() => {}).finally(() => { loading.value = false })
}

function submitSeckill() {
    if (!sform.product_id) { alert('请输入商品ID'); return }
    saving.value = true
    const payload = { ...sform }
    const promise = showEditForm.value && editingId.value
        ? adminSeckillUpdate(editingId.value, payload)
        : adminSeckillCreate(payload)
    promise.then(() => {
        alert('保存成功')
        showCreateForm.value = false
        showEditForm.value = false
        resetForm()
        loadList()
    }).catch(() => {}).finally(() => { saving.value = false })
}

function editSeckill(s) {
    editingId.value = s.id
    sform.product_id = s.product_id || 0
    sform.sku_id = s.sku_id || 0
    sform.price = s.price || 0
    sform.stock = s.stock || 0
    sform.limit_per_user = s.limit_per_user || 1
    sform.status = s.status || 1
    sform.start_time = s.start_time || ''
    sform.end_time = s.end_time || ''
    showCreateForm.value = false
    showEditForm.value = true
}

function deleteSeckill(id) {
    if (!confirm('确定删除？')) return
    adminSeckillDelete(id).then(() => {
        alert('删除成功')
        loadList()
    }).catch(() => {})
}

function cancelForm() {
    showCreateForm.value = false
    showEditForm.value = false
    editingId.value = 0
    resetForm()
}

onMounted(loadList)
</script>
