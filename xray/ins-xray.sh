#!/bin/bash
GITHUB_CMD="https://raw.githubusercontent.com/izulx1/autoscript/master/"

wget -q -O /etc/xray/vmess.json "${GITHUB_CMD}xray/vmess.json"
wget -q -O /etc/xray/vless.json "${GITHUB_CMD}xray/vless.json"
wget -q -O /etc/xray/trojan.json "${GITHUB_CMD}xray/trojan.json" 
wget -q -O /etc/xray/ss.json "${GITHUB_CMD}xray/ss.json"

wget -q -O /etc/systemd/system/vmess.service "${GITHUB_CMD}service/vmess.service"
wget -q -O /etc/systemd/system/vless.service "${GITHUB_CMD}service/vless.service"
wget -q -O /etc/systemd/system/trojan.service "${GITHUB_CMD}service/trojan.service"
wget -q -O /etc/systemd/system/ss.service "${GITHUB_CMD}service/ss.service"

rm -f ins-xray.sh >/dev/null
