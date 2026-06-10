<template>
    <AdminLayout :page-title="pageTitle">
        <div class="card">
            <div class="flex border-b">
                <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
                    :class="['flex-1 py-4 px-6 text-center font-medium transition-colors', 
                        activeTab === tab.key ? 'border-b-2 border-blue-600 text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-gray-700']">
                    {{ tab.label }}
                </button>
            </div>
            
            <div class="p-6">
                <div v-if="activeTab === 'basic'" class="space-y-4">
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label required">站点名称</label>
                            <input v-model="form.basic.site_name" type="text" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">站点 Logo</label>
                            <input v-model="form.basic.site_logo" type="text" class="form-input">
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label">站点 Favicon</label>
                            <input v-model="form.basic.site_favicon" type="text" class="form-input">
                        </div>
                        <div></div>
                    </div>
                    <div>
                        <label class="form-label">站点描述</label>
                        <textarea v-model="form.basic.site_description" class="form-input" rows="3"></textarea>
                    </div>
                    <div>
                        <label class="form-label">SEO 标题</label>
                        <input v-model="form.basic.seo_title" type="text" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">SEO 关键词</label>
                        <input v-model="form.basic.seo_keywords" type="text" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">SEO 描述</label>
                        <textarea v-model="form.basic.seo_description" class="form-input" rows="3"></textarea>
                    </div>
                    <div>
                        <label class="form-label">网站底部信息</label>
                        <textarea v-model="form.basic.footer_html" class="form-input" rows="4"></textarea>
                    </div>
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">网站关闭（维护模式）</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.basic.site_closed" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                </div>

                <div v-if="activeTab === 'payment'" class="space-y-6">
                    <div class="border rounded-lg p-4">
                        <div class="flex items-center justify-between mb-4">
                            <span class="font-medium">易支付</span>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input v-model="form.payment.epay_enabled" type="checkbox" class="sr-only peer">
                                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                            </label>
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                            <input v-model="form.payment.epay_url" type="text" class="form-input" placeholder="网关URL">
                            <input v-model="form.payment.epay_mchid" type="text" class="form-input" placeholder="商户ID">
                        </div>
                        <input v-model="form.payment.epay_key" type="text" class="form-input" placeholder="密钥">
                    </div>

                    <div class="border rounded-lg p-4">
                        <div class="flex items-center justify-between mb-4">
                            <span class="font-medium">支付宝</span>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input v-model="form.payment.alipay_enabled" type="checkbox" class="sr-only peer">
                                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                            </label>
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                            <input v-model="form.payment.alipay_appid" type="text" class="form-input" placeholder="AppID">
                            <div></div>
                        </div>
                        <textarea v-model="form.payment.alipay_private_key" class="form-input" rows="4" placeholder="商户私钥"></textarea>
                        <textarea v-model="form.payment.alipay_public_key" class="form-input" rows="4" placeholder="支付宝公钥"></textarea>
                        <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                            <span class="text-gray-700">沙箱模式</span>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input v-model="form.payment.alipay_sandbox" type="checkbox" class="sr-only peer">
                                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                            </label>
                        </div>
                    </div>

                    <div class="border rounded-lg p-4">
                        <div class="flex items-center justify-between mb-4">
                            <span class="font-medium">微信支付</span>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input v-model="form.payment.wechat_enabled" type="checkbox" class="sr-only peer">
                                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                            </label>
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                            <input v-model="form.payment.wechat_appid" type="text" class="form-input" placeholder="AppID">
                            <input v-model="form.payment.wechat_mchid" type="text" class="form-input" placeholder="商户号">
                        </div>
                        <input v-model="form.payment.wechat_key" type="text" class="form-input" placeholder="API密钥">
                    </div>

                    <div class="border rounded-lg p-4">
                        <div class="flex items-center justify-between mb-4">
                            <span class="font-medium">余额支付</span>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input v-model="form.payment.balance_enabled" type="checkbox" class="sr-only peer">
                                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                            </label>
                        </div>
                        <div>
                            <label class="form-label">充值优惠档位</label>
                            <textarea v-model="form.payment.balance_recharge_rules" class="form-input" rows="3" placeholder="每行一条，格式：充值金额:赠送金额&#10;例如：100:10"></textarea>
                        </div>
                    </div>
                </div>

                <div v-if="activeTab === 'mail'" class="space-y-4">
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label">SMTP 服务器</label>
                            <input v-model="form.mail.smtp_host" type="text" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">SMTP 端口</label>
                            <input v-model.number="form.mail.smtp_port" type="number" class="form-input">
                        </div>
                    </div>
                    <div>
                        <label class="form-label">发件人邮箱</label>
                        <input v-model="form.mail.smtp_email" type="email" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">发件人名称</label>
                        <input v-model="form.mail.smtp_name" type="text" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">密码</label>
                        <input v-model="form.mail.smtp_password" type="password" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">加密方式</label>
                        <select v-model="form.mail.smtp_encryption" class="form-input">
                            <option value="SSL">SSL</option>
                            <option value="TLS">TLS</option>
                            <option value="none">无</option>
                        </select>
                    </div>
                    <button type="button" @click="testEmail" class="btn-secondary">测试邮件发送</button>
                </div>

                <div v-if="activeTab === 'distribution'" class="space-y-4">
                    <div class="grid grid-cols-3 gap-4">
                        <div>
                            <label class="form-label">一级佣金比例 (%)</label>
                            <input v-model.number="form.distribution.distrib_level1_rate" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">二级佣金比例 (%)</label>
                            <input v-model.number="form.distribution.distrib_level2_rate" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">三级佣金比例 (%)</label>
                            <input v-model.number="form.distribution.distrib_level3_rate" type="number" class="form-input">
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label">最低提现金额</label>
                            <input v-model.number="form.distribution.distrib_min_withdraw" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">提现手续费 (%)</label>
                            <input v-model.number="form.distribution.distrib_withdraw_fee_rate" type="number" class="form-input">
                        </div>
                    </div>
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">推广员申请审核</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.distribution.distrib_audit_required" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                    <div>
                        <label class="form-label">推广链接短链域名</label>
                        <input v-model="form.distribution.distrib_short_domain" type="text" class="form-input" placeholder="如：t.example.com">
                    </div>
                </div>

                <div v-if="activeTab === 'points'" class="space-y-4">
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label">消费1元获得积分</label>
                            <input v-model.number="form.points.points_per_yuan" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">积分抵扣比例（100积分=X元）</label>
                            <input v-model.number="form.points.points_deduct_ratio" type="number" class="form-input">
                        </div>
                    </div>
                    <div>
                        <label class="form-label">最高抵扣比例 (%)</label>
                        <input v-model.number="form.points.points_max_deduct_percent" type="number" class="form-input">
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="form-label">每日签到基础积分</label>
                            <input v-model.number="form.points.signin_points" type="number" class="form-input">
                        </div>
                        <div>
                            <label class="form-label">连续签到额外奖励 (JSON)</label>
                            <input v-model="form.points.signin_continuous_reward" type="text" class="form-input">
                        </div>
                    </div>
                </div>

                <div v-if="activeTab === 'seckill'" class="space-y-4">
                    <div>
                        <label class="form-label">秒杀订单自动关闭时间（分钟）</label>
                        <input v-model.number="form.seckill.seckill_order_timeout" type="number" class="form-input">
                    </div>
                    <div>
                        <label class="form-label">秒杀并发限流（每秒请求数）</label>
                        <input v-model.number="form.seckill.seckill_rate_limit" type="number" class="form-input">
                    </div>
                </div>

                <div v-if="activeTab === 'security'" class="space-y-4">
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">登录验证码</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.security.captcha_enabled" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                    <div>
                        <label class="form-label">IP 限流阈值（每分钟请求数）</label>
                        <input v-model.number="form.security.ip_rate_limit" type="number" class="form-input">
                    </div>
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">CSRF 防护</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.security.csrf_enabled" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">敏感操作二次验证</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.security.two_factor_enabled" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                </div>

                <div v-if="activeTab === 'other'" class="space-y-4">
                    <div>
                        <label class="form-label">订单超时关闭时间（分钟）</label>
                        <input v-model.number="form.other.order_timeout_minutes" type="number" class="form-input">
                    </div>
                    <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                        <span class="text-gray-700">自动备份</span>
                        <label class="relative inline-flex items-center cursor-pointer">
                            <input v-model="form.other.auto_backup_enabled" type="checkbox" class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                        </label>
                    </div>
                    <div>
                        <label class="form-label">备份时间</label>
                        <input v-model="form.other.backup_time" type="text" class="form-input" placeholder="如：02:00">
                    </div>
                </div>

                <div v-if="activeTab === 'template'" class="space-y-4">
                    <div>
                        <label class="form-label">前台模板</label>
                        <select v-model="form.template.frontend_template" class="form-input">
                            <option value="default">默认模板</option>
                        </select>
                    </div>
                    <div>
                        <label class="form-label">后台模板</label>
                        <select v-model="form.template.backend_template" class="form-input">
                            <option value="default">默认模板</option>
                        </select>
                    </div>
                </div>

                <div class="flex justify-end mt-6">
                    <button type="button" @click="saveSettings" :disabled="loading" class="btn-primary">
                        {{ loading ? '保存中...' : '保存设置' }}
                    </button>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>import { ref, reactive, onMounted } from 'vue';
import AdminLayout from '../components/AdminLayout.vue';
import { request } from '../api/admin';
const tabs = [
 { key: 'basic', label: '基本设置' },
 { key: 'payment', label: '支付设置' },
 { key: 'mail', label: '邮件设置' },
 { key: 'distribution', label: '分销设置' },
 { key: 'points', label: '积分设置' },
 { key: 'seckill', label: '秒杀设置' },
 { key: 'security', label: '安全设置' },
 { key: 'other', label: '其他设置' },
 { key: 'template', label: '模板设置' }
];
const activeTab = ref('basic');
const loading = ref(false);
const form = reactive({
 basic: {},
 payment: {},
 mail: {},
 distribution: {},
 points: {},
 seckill: {},
 security: {},
 other: {},
 template: {}
});
const pageTitle = '系统设置';
function loadSettings(group) {
 request.get('/admin/api/settings/' + group).then(data => {
 Object.assign(form[group], data);
 }).catch(() => { });
}
function saveSettings() {
 loading.value = true;
 request.post('/admin/api/settings/' + activeTab.value, form[activeTab.value]).then(() => {
 alert('保存成功');
 }).catch(err => {
 alert(err.message || '保存失败');
 }).finally(() => {
 loading.value = false;
 });
}
function testEmail() {
 const toEmail = prompt('请输入测试邮箱地址：');
 if (!toEmail)
 return;
 request.post('/admin/api/settings/test-email', { to_email: toEmail }).then(() => {
 alert('测试邮件发送成功');
 }).catch(err => {
 alert('发送失败：' + err.message);
 });
}
activeTab.value = 'basic';
loadSettings('basic');
activeTab.value = 'payment';
loadSettings('payment');
activeTab.value = 'mail';
loadSettings('mail');
activeTab.value = 'distribution';
loadSettings('distribution');
activeTab.value = 'points';
loadSettings('points');
activeTab.value = 'seckill';
loadSettings('seckill');
activeTab.value = 'security';
loadSettings('security');
activeTab.value = 'other';
loadSettings('other');
activeTab.value = 'template';
loadSettings('template');
activeTab.value = 'basic';
</script>