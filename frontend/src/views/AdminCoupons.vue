<template>
    <AdminLayout page-title="优惠券管理">
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <h2 class="text-lg font-semibold text-gray-800">优惠券列表</h2>
                <button @click="showCreateForm = true" class="btn-primary btn-sm">新增优惠券</button>
            </div>

            <div v-if="showCreateForm || showEditForm" class="card">
                <div class="card-header font-semibold">{{ showEditForm ? '编辑优惠券' : '新增优惠券' }}</div>
                <div class="card-body space-y-3">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                        <div>
                            <label class="form-label">优惠券名称</label>
                            <input v-model="cform.name" type="text" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">类型</label>
                            <select v-model.number="cform.type" class="form-input">
                                <option :value="1">满减券</option>
                                <option :value="2">折扣券</option>
                            </select>
                        </div>
                        <div>
                            <label class="form-label">面额 / 折扣率</label>
                            <input v-model.number="cform.discount_amount" type="number" step="0.01" class="form-input" placeholder="满减券填金额，折扣券填折扣率(0-10)">
                        </div>
                        <div>
                            <label class="form-label">最低消费</label>
                            <input v-model.number="cform.min_amount" type="number" step="0.01" class="form-input" placeholder="0 表示不限制">
                        </div>
                        <div>
                            <label class="form-label">发放数量</label>
                            <input v-model.number="cform.total_count" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">有效期至</label>
                            <input v-model="cform.expire_time" type="datetime-local" class="form-input">
                        </div>
                        <div class="md:col-span-2">
                            <label class="form-label">描述</label>
                            <input v-model="cform.description" type="text" class="form-input">
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <button @click="submitCoupon" :disabled="saving" class="btn-primary btn-sm">
                            {{ saving ? '保存中...' : '保存' }}
                        </button>
                        <button @click="cancelCouponForm" class="btn-secondary btn-sm">取消</button>
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
                                <th>名称</th>
                                <th>类型</th>
                                <th>面额/折扣</th>
                                <th>最低消费</th>
                                <th>数量</th>
                                <th>已领取</th>
                                <th>状态</th>
                                <th>有效期</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="coupons.length === 0">
                                <td colspan="10" class="text-center text-gray-500 py-6">暂无数据</td>
                            </tr>
                            <tr v-for="c in coupons" :key="c.id">
                                <td>{{ c.id }}</td>
                                <td>{{ c.name }}</td>
                                <td>{{ c.type === 2 ? '折扣券' : '满减券' }}</td>
                                <td>{{ c.type === 2 ? (c.discount_rate || c.discount_amount) + '折' : '￥' + (c.discount_amount || 0) }}</td>
                                <td>￥{{ c.min_amount || 0 }}</td>
                                <td>{{ c.total_count || 0 }}</td>
                                <td>{{ c.used_count || 0 }}</td>
                                <td>{{ c.status === 1 ? '启用' : '禁用' }}</td>
                                <td class="text-xs">{{ c.expire_time || '-' }}</td>
                                <td>
                                    <button @click="editCoupon(c)" class="btn-sm btn-secondary">编辑</button>
                                    <button @click="deleteCoupon(c.id)" class="btn-sm btn-danger ml-1">删除</button>
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
import { adminCouponList, adminCouponCreate, adminCouponUpdate, adminCouponDelete } from '../api/coupon'

const coupons = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreateForm = ref(false)
const showEditForm = ref(false)
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

function submitCoupon() {
    if (!cform.name) { alert('请输入名称'); return }
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
    const promise = showEditForm.value && editingId.value
        ? adminCouponUpdate(editingId.value, payload)
        : adminCouponCreate(payload)
    promise.then(() => {
        alert('保存成功')
        showCreateForm.value = false
        showEditForm.value = false
        resetForm()
        loadList()
    }).catch(() => {}).finally(() => { saving.value = false })
}

function editCoupon(c) {
    editingId.value = c.id
    cform.name = c.name || ''
    cform.type = c.type || 1
    cform.discount_amount = c.type === 2 ? (c.discount_rate || 0) : (c.discount_amount || 0)
    cform.min_amount = c.min_amount || 0
    cform.total_count = c.total_count || 0
    cform.expire_time = c.expire_time || ''
    cform.description = c.description || ''
    showCreateForm.value = false
    showEditForm.value = true
}

function deleteCoupon(id) {
    if (!confirm('确定删除该优惠券？')) return
    adminCouponDelete(id).then(() => {
        alert('删除成功')
        loadList()
    }).catch(() => {})
}

function cancelCouponForm() {
    showCreateForm.value = false
    showEditForm.value = false
    editingId.value = 0
    resetForm()
}

onMounted(loadList)
</script>
