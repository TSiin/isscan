#!/bin/bash

# 设置变量
SERVICE_NAME="isscan"
SERVICE_PATH="/usr/local/bin/$SERVICE_NAME"
CONFIG_PATH="/etc/$SERVICE_NAME"
LOG_PATH="/var/log/$SERVICE_NAME"
SERVICE_USER="isscan"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 检查root权限
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用root权限运行此脚本${NC}"
    exit 1
fi

# 创建服务用户
echo "创建服务用户..."
id -u $SERVICE_USER >/dev/null 2>&1 || useradd -r -s /bin/false $SERVICE_USER

# 创建必要的目录
echo "创建目录..."
mkdir -p $CONFIG_PATH
mkdir -p $LOG_PATH

# 复制文件
echo "复制文件..."
cp isscan $SERVICE_PATH
cp config.yaml $CONFIG_PATH/

# 设置权限
echo "设置权限..."
chown -R $SERVICE_USER:$SERVICE_USER $CONFIG_PATH
chown -R $SERVICE_USER:$SERVICE_USER $LOG_PATH
chmod 755 $SERVICE_PATH

# 创建systemd服务
echo "创建系统服务..."
cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=Server Manager Service
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
ExecStart=$SERVICE_PATH -config $CONFIG_PATH/config.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 重新加载systemd
systemctl daemon-reload

# 启动服务
echo "启动服务..."
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

# 检查服务状态
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${GREEN}服务部署成功!${NC}"
else
    echo -e "${RED}服务启动失败，请检查日志${NC}"
    exit 1
fi

echo "
服务信息:
- 配置文件: $CONFIG_PATH/config.yaml
- 日志目录: $LOG_PATH
- 服务控制: systemctl {start|stop|restart|status} $SERVICE_NAME
" 