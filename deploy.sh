#!/bin/sh



env GOOS=linux GOARCH=arm GOARM=7 go build -o /tmp/sprinklerpi-arm7
ssh root@192.168.20.13 "systemctl stop sprinklerpi"
scp /tmp/sprinklerpi-arm7 root@192.168.20.13:/usr/local/bin/sprinklerpi
ssh root@192.168.20.13 "mkdir -p /etc/sprinklerpi"
scp config.json root@192.168.20.13:/etc/sprinklerpi/
#ssh root@192.168.20.13 "systemctl start sprinklerpi"
