#!/usr/bin/env python3
import re
import os
import json

VUE_FILES = [
    "/workspace/frontend/src/views/AdminPointsSettings.vue",
    "/workspace/frontend/src/views/Distribution.vue",
    "/workspace/frontend/src/views/UserCards.vue",
    "/workspace/frontend/src/views/Home.vue",
    "/workspace/frontend/src/views/AdminProductEdit.vue",
    "/workspace/frontend/src/views/AdminCards.vue",
    "/workspace/frontend/src/views/AdminOrders.vue",
    "/workspace/frontend/src/views/AdminCoupons.vue",
    "/workspace/frontend/src/views/Checkout.vue",
    "/workspace/frontend/src/views/ProductDetail.vue",
    "/workspace/frontend/src/views/UserCoupons.vue",
    "/workspace/frontend/src/views/AdminSettings.vue",
    "/workspace/frontend/src/views/AdminDashboard.vue",
    "/workspace/frontend/src/views/UserOrders.vue",
    "/workspace/frontend/src/views/AdminSeckills.vue",
    "/workspace/frontend/src/views/AdminLogin.vue",
    "/workspace/frontend/src/views/ProductList.vue",
    "/workspace/frontend/src/views/UserCenter.vue",
    "/workspace/frontend/src/views/CouponRedeem.vue",
    "/workspace/frontend/src/views/AdminDistribution.vue",
    "/workspace/frontend/src/views/AdminProducts.vue",
    "/workspace/frontend/src/views/AdminCategories.vue",
    "/workspace/frontend/src/views/NotFound.vue",
    "/workspace/frontend/src/views/UserPoints.vue",
    "/workspace/frontend/src/views/Cart.vue",
    "/workspace/frontend/src/views/Register.vue",
    "/workspace/frontend/src/views/Seckill.vue",
    "/workspace/frontend/src/views/Login.vue",
    "/workspace/frontend/src/views/OrderDetail.vue",
    "/workspace/frontend/src/views/UserSignin.vue",
    "/workspace/frontend/src/views/AdminUsers.vue",
    "/workspace/frontend/src/App.vue",
    "/workspace/frontend/src/components/AdminLayout.vue",
    "/workspace/frontend/src/components/Layout.vue",
]

TAILWIND_CLASSES = [
    "mb-4", "flex", "justify-between", "items-center", "gap",
    "text-blue-600", "text-red-500", "text-green-500", "text-gray",
    "font-semibold", "font-bold", "mt-2", "mt-4", "px-4", "py-2", "p-4",
    "space-y", "w-full", "w-1/2", "grid", "rounded", "border",
    "bg-white", "bg-gray", "hover:"
]

def read_file(path):
    with open(path, "r", encoding="utf-8") as f:
        return f.read()

def split_sections(content):
    template_match = re.search(r"<template[^>]*>(.*?)</template>", content, re.DOTALL)
    script_match = re.search(r"<script[^>]*>(.*?)</script>", content, re.DOTALL)
    template = template_match.group(1) if template_match else ""
    script = script_match.group(1) if script_match else ""
    return template, script

def get_line_number(content, pattern_start, search_from=0):
    return content[:pattern_start].count("\n") + 1

# ============ 检查 1：Tailwind 残留类 ============
def check_tailwind(filepath, content, template):
    issues = []
    # 在模板中查找 class 属性
    class_patterns = re.finditer(r'class\s*=\s*["\']([^"\']+)["\']', template)
    for m in class_patterns:
        class_value = m.group(1)
        line_no = get_line_number(content, m.start(1))
        classes = class_value.split()
        for cls in classes:
            cls_clean = cls.strip()
            for tw in TAILWIND_CLASSES:
                if cls_clean == tw or cls_clean.startswith(tw + "-") or cls_clean.startswith(tw):
                    # 精确匹配
                    if cls_clean == tw or (tw in ["gap", "space-y", "text-gray", "bg-gray", "rounded", "border"] and cls_clean.startswith(tw)) or (tw == "hover:" and cls_clean.startswith("hover:")) or (tw == "grid" and cls_clean == "grid") or (tw == "flex" and cls_clean == "flex"):
                        if any(cls_clean == x or x in cls_clean for x in TAILWIND_CLASSES):
                            pass
                        issues.append({
                            "file": filepath,
                            "line": line_no,
                            "type": "Tailwind残留",
                            "desc": f"class='{class_value}' 中包含 Tailwind 类 '{cls_clean}'（无样式效果）"
                        })
                        break
    # 去重
    seen = set()
    unique_issues = []
    for issue in issues:
        key = (issue["file"], issue["line"], issue["desc"])
        if key not in seen:
            seen.add(key)
            unique_issues.append(issue)
    return unique_issues

# ============ 检查 2：功能性绑定 ============
def check_functional(filepath, content, template, script):
    issues = []

    # 2a: @selection-change
    for m in re.finditer(r'@selection-change\s*=\s*["\']([^"\']+)["\']', template):
        handler = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        # 提取函数名（去掉括号和参数）
        func_name = re.sub(r'\(.*$', '', handler).strip()
        if func_name and not re.search(r'def\s+' + re.escape(func_name) + r'\b|const\s+' + re.escape(func_name) + r'\s*=|function\s+' + re.escape(func_name) + r'\b|\b' + re.escape(func_name) + r'\s*=\s*\(|' + re.escape(func_name) + r'\s*\([^)]*\)\s*\{', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"@selection-change=\"{handler}\" 绑定的方法 '{func_name}' 在 script 中未定义"
            })

    # 2b: @size-change
    for m in re.finditer(r'@size-change\s*=\s*["\']([^"\']+)["\']', template):
        handler = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        func_name = re.sub(r'\(.*$', '', handler).strip()
        if func_name and not re.search(r'\b' + re.escape(func_name) + r'\s*[=(]', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"@size-change=\"{handler}\" 绑定的方法 '{func_name}' 在 script 中未定义"
            })

    # 2c: @current-change
    for m in re.finditer(r'@current-change\s*=\s*["\']([^"\']+)["\']', template):
        handler = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        func_name = re.sub(r'\(.*$', '', handler).strip()
        if func_name and not re.search(r'\b' + re.escape(func_name) + r'\s*[=(]', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"@current-change=\"{handler}\" 绑定的方法 '{func_name}' 在 script 中未定义"
            })

    # 2d: v-model:current-page / v-model:page-size
    for m in re.finditer(r'v-model:current-page\s*=\s*["\']([^"\']+)["\']', template):
        var = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        if not re.search(r'\b' + re.escape(var) + r'\b', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"v-model:current-page=\"{var}\" 变量 '{var}' 在 script 中未定义"
            })

    for m in re.finditer(r'v-model:page-size\s*=\s*["\']([^"\']+)["\']', template):
        var = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        if not re.search(r'\b' + re.escape(var) + r'\b', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"v-model:page-size=\"{var}\" 变量 '{var}' 在 script 中未定义"
            })

    # 2e: el-form rules
    for m in re.finditer(r'<el-form[^>]*\brules\s*=\s*["\']([^"\']+)["\'][^>]*>', template):
        var = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        if not re.search(r'\b' + re.escape(var) + r'\b', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"el-form :rules=\"{var}\" 规则对象 '{var}' 在 script 中未定义"
            })

    # 2f: el-dialog v-model
    for m in re.finditer(r'<el-dialog[^>]*\bv-model\s*=\s*["\']([^"\']+)["\'][^>]*>', template):
        var = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        if not re.search(r'\b' + re.escape(var) + r'\b', script):
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "功能性错误",
                "desc": f"el-dialog v-model=\"{var}\" 变量 '{var}' 在 script 中未定义"
            })

    return issues

# ============ 检查 3：模板变量是否在 script 中定义 ============
def check_variables(filepath, content, template, script):
    issues = []

    # 提取 script 中定义的所有变量名
    defined_vars = set()

    # ref() / reactive() / computed() 定义
    for m in re.finditer(r'(?:const|let|var)\s+(\w+)\s*=\s*(?:ref|reactive|computed)\s*\(', script):
        defined_vars.add(m.group(1))

    # 普通 const/let/var 定义
    for m in re.finditer(r'(?:const|let|var)\s+(\w+)\s*=', script):
        defined_vars.add(m.group(1))

    # function 定义
    for m in re.finditer(r'(?:async\s+)?function\s+(\w+)\s*\(', script):
        defined_vars.add(m.group(1))

    # 箭头函数定义 const foo = () => {} 或 const foo = async () => {}
    for m in re.finditer(r'const\s+(\w+)\s*=\s*(?:async\s*)?\(', script):
        defined_vars.add(m.group(1))

    # import 引入的
    for m in re.finditer(r'import\s+(?:{[^}]*}\s+from\s+|)([\w]+)\s+from', script):
        defined_vars.add(m.group(1))
    for m in re.finditer(r'import\s*{\s*([^}]+)\s*}\s+from', script):
        for name in m.group(1).split(','):
            name = name.strip()
            if name:
                defined_vars.add(name.split(' as ')[-1].strip())

    # defineProps / defineEmits
    props_match = re.search(r'defineProps\s*<\s*\{([^}]+)\}', script)
    if props_match:
        for m in re.finditer(r'(\w+)\s*[?:]', props_match.group(1)):
            defined_vars.add(m.group(1))
    # defineProps({ ... }) 对象式
    props_obj_match = re.search(r'defineProps\s*\(\s*\{([^}]+)\}', script)
    if props_obj_match:
        for m in re.finditer(r'(\w+)\s*:', props_obj_match.group(1)):
            defined_vars.add(m.group(1))

    # 提取模板中使用的变量（v-model, v-if, v-for, {{ }}, :xxx="yyy")
    template_vars = set()

    # v-model
    for m in re.finditer(r'v-model(?::\w+)?\s*=\s*["\']([^"\']+)["\']', template):
        expr = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        # 提取顶层变量名
        var = re.split(r'[.\[(!]', expr)[0].strip()
        if var and not re.match(r'^(true|false|null|undefined|this)$', var) and var not in defined_vars:
            template_vars.add(var)
            issues.append({
                "file": filepath,
                "line": line_no,
                "type": "变量未定义",
                "desc": f"模板 v-model 使用变量 '{var}' (表达式: {expr})，但在 script setup 中未用 ref/computed 定义"
            })

    # v-if / v-show
    for m in re.finditer(r'v-(?:if|show)\s*=\s*["\']([^"\']+)["\']', template):
        expr = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        var = re.split(r'[.\[(!=<>]', expr)[0].strip()
        if var and not re.match(r'^(true|false|null|undefined|this|!)$', var) and var not in defined_vars and not var.startswith('!'):
            # 去掉 ! 前缀
            actual_var = var.lstrip('!')
            actual_var = re.split(r'[.\[(]', actual_var)[0].strip()
            if actual_var and actual_var not in defined_vars and not re.match(r'^(true|false|null|undefined)$', actual_var):
                issues.append({
                    "file": filepath,
                    "line": line_no,
                    "type": "变量未定义",
                    "desc": f"模板 v-if/v-show 使用变量 '{actual_var}' (表达式: {expr})，但在 script setup 中未定义"
                })

    # v-for
    for m in re.finditer(r'v-for\s*=\s*["\']([^"\']+)["\']', template):
        expr = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        # v-for="item in items"  - 取 in 后面的
        in_match = re.search(r'\s+(?:in|of)\s+(.+)$', expr)
        if in_match:
            iter_expr = in_match.group(1).strip()
            var = re.split(r'[.\[(!]', iter_expr)[0].strip()
            if var and var not in defined_vars and not re.match(r'^(true|false|null|undefined)$', var):
                issues.append({
                    "file": filepath,
                    "line": line_no,
                    "type": "变量未定义",
                    "desc": f"模板 v-for 迭代变量 '{var}' (表达式: {iter_expr})，但在 script setup 中未定义"
                })

    # 属性绑定 :prop="xxx"
    for m in re.finditer(r'(?<!@)\s+:[\w-]+\s*=\s*["\']([^"\']+)["\']', template):
        expr = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        # 跳过纯字符串字面量 / 数字 / 事件绑定
        if expr.startswith("'") or expr.startswith('"') or re.match(r'^\d+$', expr):
            continue
        var = re.split(r'[.\[(!=<>+\-*/%]', expr)[0].strip()
        var = var.lstrip('!')
        if var and var not in defined_vars and not re.match(r'^(true|false|null|undefined|this)$', var) and len(var) > 0:
            # 只报告第一个发现
            pass

    # 双花括号 {{ xxx }}
    for m in re.finditer(r'\{\{\s*([^}]+?)\s*\}\}', template):
        expr = m.group(1).strip()
        line_no = get_line_number(content, m.start(1))
        var = re.split(r'[.\[(!=<>+\-*/%]', expr)[0].strip()
        var = var.lstrip('!')
        if var and var not in defined_vars and not re.match(r'^(true|false|null|undefined|this|\d)$', var) and len(var) > 0 and not re.match(r'^[\W_]+$', var):
            # 去重检查在主循环外
            pass

    return issues

# ============ 主检查循环 ============
all_issues = []
for filepath in VUE_FILES:
    if not os.path.exists(filepath):
        print(f"SKIP: {filepath}")
        continue
    content = read_file(filepath)
    template, script = split_sections(content)

    issues_1 = check_tailwind(filepath, content, template)
    issues_2 = check_functional(filepath, content, template, script)
    issues_3 = check_variables(filepath, content, template, script)

    all_issues.extend(issues_1)
    all_issues.extend(issues_2)
    all_issues.extend(issues_3)

# 去重变量未定义问题
seen_vars = set()
final_issues = []
for issue in all_issues:
    key = (issue["file"], issue["type"], issue["desc"])
    if key not in seen_vars:
        seen_vars.add(key)
        final_issues.append(issue)

# 按文件分组输出
from collections import defaultdict
by_file = defaultdict(list)
for issue in final_issues:
    by_file[issue["file"]].append(issue)

print("=" * 80)
print("VUE 文件系统性检查报告")
print(f"扫描文件数: {len(VUE_FILES)}")
print(f"发现问题数: {len(final_issues)}")
print("=" * 80)

# 核心页面优先
CORE_FILES = [
    "AdminProducts.vue", "AdminOrders.vue", "AdminUsers.vue",
    "AdminCoupons.vue", "AdminSeckills.vue", "UserOrders.vue", "ProductList.vue"
]

# 先输出核心页面
print("\n### 核心页面问题汇总 ###\n")
for core in CORE_FILES:
    core_path = None
    for f in by_file.keys():
        if f.endswith(core):
            core_path = f
            break
    if core_path and core_path in by_file:
        print(f"\n{'='*60}")
        print(f"★ {core_path}")
        print(f"{'='*60}")
        for issue in sorted(by_file[core_path], key=lambda x: x["line"]):
            print(f"  [第{issue['line']:3d}行] [{issue['type']}] {issue['desc']}")

print("\n\n### 其他文件问题汇总 ###\n")
for filepath in sorted(by_file.keys()):
    is_core = any(filepath.endswith(c) for c in CORE_FILES)
    if is_core:
        continue
    print(f"\n--- {filepath} ({len(by_file[filepath])} 个问题) ---")
    for issue in sorted(by_file[filepath], key=lambda x: x["line"]):
        print(f"  [第{issue['line']:3d}行] [{issue['type']}] {issue['desc']}")

# 统计
print("\n" + "=" * 80)
print("统计:")
type_counts = defaultdict(int)
for issue in final_issues:
    type_counts[issue["type"]] += 1
for t, c in sorted(type_counts.items(), key=lambda x: -x[1]):
    print(f"  {t}: {c} 处")
print("=" * 80)

# 输出为 JSON 以便进一步处理
with open("/workspace/vue_check_report.json", "w", encoding="utf-8") as f:
    json.dump(final_issues, f, ensure_ascii=False, indent=2)
print(f"\n详细 JSON 报告已保存到: /workspace/vue_check_report.json")
