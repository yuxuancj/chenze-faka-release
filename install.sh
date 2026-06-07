#!/bin/bash
echo "晨泽发卡系统 - 安装脚本"
echo ""
read -p "请输入数据库主机 [localhost]: " DB_HOST
DB_HOST=${DB_HOST:-localhost}
read -p "请输入数据库端口 [3306]: " DB_PORT
DB_PORT=${DB_PORT:-3306}
read -p "请输入数据库用户名 [root]: " DB_USER
DB_USER=${DB_USER:-root}
read -s -p "请输入数据库密码: " DB_PASS
echo ""
read -p "请输入数据库名 [chenze_faka]: " DB_NAME
DB_NAME=${DB_NAME:-chenze_faka}

# 更新配置文件
sed -i "s/host: localhost/host: $DB_HOST/" config.yaml
sed -i "s/port: 3306/port: $DB_PORT/" config.yaml
sed -i "s/user: root/user: $DB_USER/" config.yaml
sed -i "s/password: /password: $DB_PASS/" config.yaml
sed -i "s/dbname: chenze_faka/dbname: $DB_NAME/" config.yaml

echo ""
echo "配置已更新！请运行 ./chenze-faka 启动系统"
