sudo systemctl disable stun-server.service
sudo systemctl stop stun-server.service

sudo mv stun-server /usr/local/bin/
sudo chmod +x /usr/local/bin/stun-server

cat > /etc/systemd/system/stun-server.service <<EOF
[Unit]
Description=STUN Server Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/stun-server -port 3478

# 【核心设置】修改停止信号为 SIGKILL（即 kill -9）
KillSignal=SIGKILL
# 【可选】确保强杀所有由该服务衍生出来的子进程
KillMode=control-group
# always on-failure
Restart=on-failure
RestartSec=1
RestartForceExitStatus=SIGINT

StandardOutput=journal
StandardError=journal
# 安全增强配置
ProtectSystem=full
PrivateTmp=true

[Install]
WantedBy=multi-user.target

EOF


sudo systemctl daemon-reload
sudo systemctl start stun-server
sudo systemctl enable stun-server

sudo systemctl status stun-server