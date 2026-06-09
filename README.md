# 晨泽发卡 - 生产部署包

## 快速启动

```bash
# 1. 修改配置
vim full/config.yaml

# 2. 启动
cd full
chmod +x chenze_faka
./chenze_faka
```

## 运行环境

- Linux amd64（已静态编译，无需 CGO）
- MySQL 5.7+

## 首次启动

1. 修改 `full/config.yaml` 中的数据库连接信息
2. 运行后访问 `http://your-server:8080/admin/` 用默认管理员登录：
   - 邮箱：`admin@chenze.com`
   - 密码：`admin123`

## systemd 部署（可选）

```ini
[Unit]
Description=Chenze Faka
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/chenze-faka
ExecStart=/opt/chenze-faka/chenze_faka
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo cp chenze_faka /usr/local/bin/
sudo cp config.yaml /etc/chenze-faka/config.yaml
sudo systemctl daemon-reload
sudo systemctl enable chenze-faka
sudo systemctl start chenze-faka
```
