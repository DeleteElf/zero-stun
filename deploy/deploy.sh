systemctl disable stun-server.service
systemctl stop stun-server.service
DIR=/usr/local/bin
rm -f ${DIR}/stun-server
cp stun-server ${DIR}/
chmod +x ${DIR}/stun-server
cp -n ../etc/server.yaml ${DIR}/etc/

cat > /etc/systemd/system/stun-server.service <<EOF
[Unit]
Description=STUN Server Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=${DIR}/stun-server
WorkingDirectory=${DIR}/
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


systemctl daemon-reload
systemctl start stun-server
systemctl enable stun-server

systemctl status stun-server