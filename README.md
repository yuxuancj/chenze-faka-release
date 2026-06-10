# 晨泽发卡系统发布包

## 发布包说明

- **chenze_faka_v2.3-license_prod.tar.gz** (约 6.5MB)
  - 版本：v2.3-license
  - 类型：静态二进制（CGO_ENABLED=0，纯 Go，无 glibc 依赖）
  - 系统：Linux x86_64，CentOS 7/8/9 / Ubuntu / 宝塔面板

## 新增功能

1. **授权码验证**：安装向导增加授权码字段
2. **在线更新**：后台 /admin/api/update/ 接口
3. **完整数据库迁移**：19 张表
4. **静态编译**：不依赖系统 glibc

## 使用方法

```
