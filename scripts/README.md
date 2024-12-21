
cp ./issue2md.service /etc/systemd/system/
systemctl enable issue2md.service
systemctl daemon-reload
systemctl start issue2md

