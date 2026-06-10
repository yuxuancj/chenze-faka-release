<template>
    <AdminLayout :page-title="isEdit ? '编辑商品' : '新增商品'">
        <div class="card max-w-4xl">
            <div class="card-body">
                <div v-if="loading" class="text-center py-12 text-gray-500">
                    <div class="animate-spin h-12 w-12 mx-auto"></div>
                    <p class="mt-4">加载中...</p>
                </div>
                <form v-else @submit.prevent="save" class="space-y-6">
                    <div class="border rounded-lg overflow-hidden">
                        <button type="button" @click="expandedSections.basic = !expandedSections.basic" 
                            class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors">
                            <span class="font-medium">基本信息</span>
                            <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': expandedSections.basic }">▼</span>
                        </button>
                        <div v-show="expandedSections.basic" class="p-6 space-y-4 border-t">
                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="form-label required">商品名称</label>
                                    <input v-model="form.name" type="text" class="form-input" placeholder="请输入商品名称">
                                </div>
                                <div>
                                    <label class="form-label required">商品分类</label>
                                    <select v-model.number="form.category_id" class="form-input">
                                        <option value="0">请选择分类</option>
                                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                                            {{ cat.name }}
                                        </option>
                                    </select>
                                </div>
                            </div>
                            <div class="grid grid-cols-3 gap-4">
                                <div>
                                    <label class="form-label required">价格</label>
                                    <input v-model.number="form.price" type="number" step="0.01" class="form-input" placeholder="0.00">
                                </div>
                                <div>
                                    <label class="form-label required">库存</label>
                                    <input v-model.number="form.stock" type="number" class="form-input" placeholder="0">
                                </div>
                                <div>
                                    <label class="form-label">排序</label>
                                    <input v-model.number="form.sort" type="number" class="form-input" placeholder="0" title="数字越小越靠前">
                                </div>
                            </div>
                            <div>
                                <label class="form-label">商品图片</label>
                                <div class="flex items-center gap-3">
                                    <input v-model="form.image" type="text" class="form-input flex-1" placeholder="图片URL或上传">
                                    <button type="button" class="btn-secondary">上传</button>
                                </div>
                            </div>
                            <div>
                                <label class="form-label">商品描述</label>
                                <textarea v-model="form.description" class="form-input" rows="4" placeholder="请输入商品描述"></textarea>
                            </div>
                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="form-label">状态</label>
                                    <select v-model.number="form.status" class="form-input">
                                        <option :value="1">上架</option>
                                        <option :value="0">下架</option>
                                    </select>
                                </div>
                                <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                                    <span class="text-gray-700">隐藏商品</span>
                                    <label class="relative inline-flex items-center cursor-pointer">
                                        <input v-model="form.is_hidden" type="checkbox" class="sr-only peer">
                                        <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="border rounded-lg overflow-hidden">
                        <button type="button" @click="expandedSections.delivery = !expandedSections.delivery" 
                            class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors">
                            <span class="font-medium">发货类型</span>
                            <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': expandedSections.delivery }">▼</span>
                        </button>
                        <div v-show="expandedSections.delivery" class="p-6 border-t">
                            <div class="flex border-b mb-4">
                                <button type="button" v-for="dt in deliveryTypes" :key="dt.value"
                                    @click="form.delivery_type = dt.value"
                                    :class="['flex-1 py-3 px-4 text-center font-medium transition-colors', 
                                        form.delivery_type === dt.value ? 'border-b-2 border-blue-600 text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-gray-700']">
                                    {{ dt.label }}
                                </button>
                            </div>

                            <div v-if="form.delivery_type === 'card'" class="mt-4 space-y-4">
                                <div>
                                    <label class="form-label">卡密批量导入</label>
                                    <div class="border border-dashed rounded-lg p-6 text-center">
                                        <input type="file" accept=".txt" @change="handleCardFileUpload" class="hidden" id="card-file">
                                        <label for="card-file" class="cursor-pointer">
                                            <div class="text-4xl mb-2">📁</div>
                                            <p class="text-gray-600">点击或拖拽上传 .txt 文件</p>
                                            <p class="text-sm text-gray-400 mt-1">每行一个卡密，或使用 | 分隔卡号和密码</p>
                                        </label>
                                    </div>
                                    <div v-if="cardList.length > 0" class="mt-4">
                                        <div class="flex items-center justify-between mb-2">
                                            <span class="text-sm text-gray-600">已导入 {{ cardList.length }} 条卡密</span>
                                            <button type="button" @click="cardList = []" class="text-sm text-red-500">清空</button>
                                        </div>
                                        <div class="max-h-40 overflow-y-auto border rounded-lg p-2">
                                            <div v-for="(card, index) in cardList.slice(0, 20)" :key="index" 
                                                class="text-sm text-gray-600 py-1 border-b border-gray-100 last:border-0">
                                                {{ card }}
                                            </div>
                                            <div v-if="cardList.length > 20" class="text-center text-gray-400 py-2">
                                                ... 还有 {{ cardList.length - 20 }} 条
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div v-if="form.delivery_type === 'text'" class="mt-4">
                                <label class="form-label">文本内容</label>
                                <textarea v-model="deliveryConfig.text_content" class="form-input" rows="8" placeholder="购买后显示的固定文本内容"></textarea>
                            </div>

                            <div v-if="form.delivery_type === 'api'" class="mt-4 space-y-4">
                                <div>
                                    <label class="form-label required">API URL</label>
                                    <input v-model="deliveryConfig.api_url" type="text" class="form-input" placeholder="https://api.example.com/auth">
                                </div>
                                <div>
                                    <label class="form-label">请求方式</label>
                                    <select v-model="deliveryConfig.api_method" class="form-input">
                                        <option value="GET">GET</option>
                                        <option value="POST">POST</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="form-label">Header 参数</label>
                                    <div class="border rounded-lg p-3">
                                        <div v-for="(value, key) in deliveryConfig.api_headers" :key="key" class="flex gap-2 mb-2">
                                            <input :value="key" @input="updateApiHeaderKey($event, key)" type="text" class="form-input flex-1" placeholder="Key">
                                            <input v-model="deliveryConfig.api_headers[key]" type="text" class="form-input flex-1" placeholder="Value">
                                            <button type="button" @click="deleteApiHeader(key)" class="btn-danger">删除</button>
                                        </div>
                                        <button type="button" @click="addApiHeader" class="btn-secondary w-full">添加 Header</button>
                                    </div>
                                </div>
                                <div>
                                    <label class="form-label">Body 参数</label>
                                    <div class="border rounded-lg p-3">
                                        <div v-for="(value, key) in deliveryConfig.api_body" :key="key" class="flex gap-2 mb-2">
                                            <input :value="key" @input="updateApiBodyKey($event, key)" type="text" class="form-input flex-1" placeholder="Key">
                                            <input v-model="deliveryConfig.api_body[key]" type="text" class="form-input flex-1" placeholder="Value">
                                            <button type="button" @click="deleteApiBody(key)" class="btn-danger">删除</button>
                                        </div>
                                        <button type="button" @click="addApiBody" class="btn-secondary w-full">添加 Body 参数</button>
                                    </div>
                                </div>
                                <div class="grid grid-cols-2 gap-4">
                                    <div>
                                        <label class="form-label">成功标识字段</label>
                                        <input v-model="deliveryConfig.api_success_field" type="text" class="form-input" placeholder="code==0">
                                    </div>
                                    <div>
                                        <label class="form-label">授权码返回字段</label>
                                        <input v-model="deliveryConfig.api_auth_field" type="text" class="form-input" placeholder="data.auth_code">
                                    </div>
                                </div>
                                <div>
                                    <label class="form-label">授权到期天数</label>
                                    <input v-model.number="deliveryConfig.api_expire_days" type="number" class="form-input" placeholder="365">
                                </div>
                            </div>

                            <div v-if="form.delivery_type === 'infinite'" class="mt-4">
                                <label class="form-label">固定内容</label>
                                <textarea v-model="deliveryConfig.text_content" class="form-input" rows="4" placeholder="用户购买后显示的固定内容"></textarea>
                            </div>

                            <div v-if="form.delivery_type === 'manual'" class="mt-4 p-4 bg-blue-50 rounded-lg">
                                <p class="text-blue-700">人工核销模式：用户购买后生成工单，管理员需手动完成核销。</p>
                            </div>
                        </div>
                    </div>

                    <div class="border rounded-lg overflow-hidden">
                        <button type="button" @click="expandedSections.sku = !expandedSections.sku" 
                            class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors">
                            <span class="font-medium">SKU 多规格管理</span>
                            <span class="flex items-center gap-2">
                                <label class="relative inline-flex items-center cursor-pointer">
                                    <input v-model="form.has_sku" type="checkbox" class="sr-only peer">
                                    <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                                </label>
                                <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': expandedSections.sku }">▼</span>
                            </span>
                        </button>
                        <div v-show="expandedSections.sku && form.has_sku" class="p-6 border-t">
                            <div class="flex justify-between items-center mb-4">
                                <div class="flex gap-2">
                                    <button type="button" @click="addSku" class="btn-primary">添加规格</button>
                                    <button type="button" @click="importSku" class="btn-secondary">批量导入</button>
                                </div>
                            </div>
                            <div class="overflow-x-auto">
                                <table class="w-full border">
                                    <thead>
                                        <tr class="bg-gray-50">
                                            <th class="border px-4 py-2 text-left">规格组合</th>
                                            <th class="border px-4 py-2 text-left">SKU编码</th>
                                            <th class="border px-4 py-2 text-left">价格</th>
                                            <th class="border px-4 py-2 text-left">库存</th>
                                            <th class="border px-4 py-2 text-left">图片</th>
                                            <th class="border px-4 py-2 text-center">操作</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr v-for="(sku, index) in form.skus" :key="index">
                                            <td class="border px-4 py-2">
                                                <input v-model="sku.spec_names" type="text" class="form-input" placeholder="如：颜色:红色,尺寸:L">
                                            </td>
                                            <td class="border px-4 py-2">
                                                <input v-model="sku.sku_code" type="text" class="form-input" placeholder="SKU编码">
                                            </td>
                                            <td class="border px-4 py-2">
                                                <input v-model.number="sku.price" type="number" step="0.01" class="form-input" placeholder="0.00">
                                            </td>
                                            <td class="border px-4 py-2">
                                                <input v-model.number="sku.stock" type="number" class="form-input" placeholder="0">
                                            </td>
                                            <td class="border px-4 py-2">
                                                <input v-model="sku.image" type="text" class="form-input" placeholder="图片URL">
                                            </td>
                                            <td class="border px-4 py-2 text-center">
                                                <button type="button" @click="removeSku(index)" class="text-red-500 hover:text-red-700">删除</button>
                                            </td>
                                        </tr>
                                        <tr v-if="form.skus.length === 0">
                                            <td colspan="6" class="border px-4 py-8 text-center text-gray-400">
                                                暂无 SKU，点击上方按钮添加
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>

                    <div class="border rounded-lg overflow-hidden">
                        <button type="button" @click="expandedSections.member = !expandedSections.member" 
                            class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors">
                            <span class="font-medium">会员等级价</span>
                            <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': expandedSections.member }">▼</span>
                        </button>
                        <div v-show="expandedSections.member" class="p-6 border-t">
                            <div class="space-y-3">
                                <div v-for="level in memberLevels" :key="level.id" class="flex items-center gap-4">
                                    <span class="w-24 text-gray-600">{{ level.name }}</span>
                                    <input :value="getMemberPrice(level.id)" @input="updateMemberPrice(level.id, $event)" type="number" step="0.01" 
                                        class="form-input w-32" :placeholder="'默认 ' + form.price">
                                </div>
                            </div>
                            <p class="text-sm text-gray-400 mt-4">留空则使用默认价格</p>
                        </div>
                    </div>

                    <div class="border rounded-lg overflow-hidden">
                        <button type="button" @click="expandedSections.seo = !expandedSections.seo" 
                            class="w-full px-6 py-4 flex items-center justify-between bg-gray-50 hover:bg-gray-100 transition-colors">
                            <span class="font-medium">SEO 设置</span>
                            <span class="text-gray-400 transition-transform" :class="{ 'rotate-180': expandedSections.seo }">▼</span>
                        </button>
                        <div v-show="expandedSections.seo" class="p-6 border-t space-y-4">
                            <div>
                                <label class="form-label">SEO 标题</label>
                                <input v-model="form.seo.title" type="text" class="form-input" placeholder="页面标题">
                            </div>
                            <div>
                                <label class="form-label">SEO 关键词</label>
                                <input v-model="form.seo.keywords" type="text" class="form-input" placeholder="多个关键词用逗号分隔">
                            </div>
                            <div>
                                <label class="form-label">SEO 描述</label>
                                <textarea v-model="form.seo.description" class="form-input" rows="3" placeholder="页面描述"></textarea>
                            </div>
                        </div>
                    </div>

                    <div class="flex items-center gap-3 pt-4">
                        <button type="submit" :disabled="loading" class="btn-primary">
                            {{ loading ? '保存中...' : '保存' }}
                        </button>
                        <router-link to="/admin/products" class="btn-secondary">返回</router-link>
                    </div>
                </form>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AdminLayout from '../components/AdminLayout.vue';
import { adminProductCreate, adminProductUpdate, adminProductDetail } from '../api/admin';
const route = useRoute();
const router = useRouter();
const isEdit = computed(() => !!route.params.id);
const loading = ref(false);
const categories = ref([]);
const memberLevels = ref([]);
const cardList = ref([]);
const expandedSections = reactive({
 basic: true,
 delivery: true,
 sku: true,
 member: true,
 seo: true
});
const deliveryTypes = [
 { value: 'card', label: '卡密' },
 { value: 'text', label: '文本' },
 { value: 'api', label: 'API授权' },
 { value: 'infinite', label: '无限卡密' },
 { value: 'manual', label: '人工核销' }
];
const deliveryConfig = reactive({
 card_list: [],
 text_content: '',
 api_url: '',
 api_method: 'POST',
 api_headers: {},
 api_body: {},
 api_success_field: '',
 api_auth_field: '',
 api_expire_days: 365
});
const form = reactive({
 name: '',
 category_id: 0,
 price: 0,
 stock: 0,
 description: '',
 image: '',
 images: [],
 status: 1,
 is_hidden: false,
 has_sku: false,
 sort: 0,
 delivery_type: 'card',
 delivery_config: deliveryConfig,
 skus: [],
 member_prices: [],
 seo: {
 title: '',
 keywords: '',
 description: ''
 }
});
function loadCategories() {
 fetch('/admin/api/categories').then(res => res.json()).then(data => {
 categories.value = data || [];
 }).catch(() => { });
}
function loadMemberLevels() {
 fetch('/admin/api/settings/member_levels').then(res => res.json()).then(data => {
 memberLevels.value = data || [];
 }).catch(() => {
 memberLevels.value = [
 { id: 1, name: '普通会员' },
 { id: 2, name: 'VIP会员' },
 { id: 3, name: '钻石会员' }
 ];
 });
}
function loadProduct() {
 if (!route.params.id)
 return;
 loading.value = true;
 adminProductDetail(route.params.id).then((data) => {
 if (data) {
 form.name = data.name || '';
 form.category_id = data.category_id || 0;
 form.price = data.price || 0;
 form.stock = data.stock || 0;
 form.description = data.description || '';
 form.image = data.image || '';
 form.status = data.status === undefined ? 1 : data.status;
 form.is_hidden = data.is_hidden || false;
 form.has_sku = data.has_sku || false;
 form.sort = data.sort || 0;
 form.delivery_type = data.delivery_type || 'card';
 if (data.delivery_config) {
 Object.assign(deliveryConfig, data.delivery_config);
 }
 if (data.skus) {
 form.skus = data.skus.map(s => ({
 spec_names: typeof s.spec_names === 'object' ? JSON.stringify(s.spec_names) : s.spec_names || '',
 sku_code: s.sku_code || '',
 price: s.price || 0,
 stock: s.stock || 0,
 image: s.image || ''
 }));
 }
 if (data.member_prices) {
 form.member_prices = data.member_prices;
 }
 if (data.seo) {
 Object.assign(form.seo, data.seo);
 }
 }
 }).catch(() => { }).finally(() => {
 loading.value = false;
 });
}
function getMemberPrice(levelId) {
 const mp = form.member_prices.find(p => p.level_id === levelId);
 return mp ? mp.price : '';
}
function updateMemberPrice(levelId, event) {
 const value = parseFloat(event.target.value);
 const existing = form.member_prices.find(p => p.level_id === levelId);
 if (value && value > 0) {
 if (existing) {
 existing.price = value;
 } else {
 form.member_prices.push({ level_id: levelId, price: value });
 }
 } else if (existing) {
 form.member_prices = form.member_prices.filter(p => p.level_id !== levelId);
 }
}
function handleCardFileUpload(event) {
 const file = event.target.files[0];
 if (!file)
 return;
 const reader = new FileReader();
 reader.onload = (e) => {
 const content = e.target.result;
 cardList.value = content.split('\n').map(line => line.trim()).filter(line => line);
 deliveryConfig.card_list = cardList.value;
 };
 reader.readAsText(file);
}
function addApiHeader() {
 deliveryConfig.api_headers['new_key'] = '';
}
function deleteApiHeader(key) {
 delete deliveryConfig.api_headers[key];
}
function updateApiHeaderKey(event, oldKey) {
 const newValue = event.target.value;
 if (newValue !== oldKey) {
 deliveryConfig.api_headers[newValue] = deliveryConfig.api_headers[oldKey];
 delete deliveryConfig.api_headers[oldKey];
 }
}
function addApiBody() {
 deliveryConfig.api_body['new_key'] = '';
}
function deleteApiBody(key) {
 delete deliveryConfig.api_body[key];
}
function updateApiBodyKey(event, oldKey) {
 const newValue = event.target.value;
 if (newValue !== oldKey) {
 deliveryConfig.api_body[newValue] = deliveryConfig.api_body[oldKey];
 delete deliveryConfig.api_body[oldKey];
 }
}
function addSku() {
 form.skus.push({
 spec_names: '',
 sku_code: '',
 price: 0,
 stock: 0,
 image: ''
 });
}
function removeSku(index) {
 form.skus.splice(index, 1);
}
function importSku() {
 alert('批量导入功能开发中，请手动添加');
}
function save() {
 if (!form.name) {
 alert('请输入商品名称');
 return;
 }
 if (form.price <= 0) {
 alert('请输入有效的价格');
 return;
 }
 if (form.category_id === 0) {
 alert('请选择商品分类');
 return;
 }
 form.delivery_config = deliveryConfig;
 const validMemberPrices = [];
 memberLevels.value.forEach(level => {
 const price = getMemberPrice(level.id);
 if (price && price > 0) {
 validMemberPrices.push({ level_id: level.id, price });
 }
 });
 form.member_prices = validMemberPrices;
 loading.value = true;
 const action = isEdit.value
 ? adminProductUpdate(route.params.id, form)
 : adminProductCreate(form);
 action.then(() => {
 alert('保存成功');
 router.push('/admin/products');
 }).catch((err) => {
 alert(err.message || '保存失败');
 }).finally(() => {
 loading.value = false;
 });
}
onMounted(() => {
 loadCategories();
 loadMemberLevels();
 loadProduct();
});
</script>