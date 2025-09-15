#!/bin/bash

go mod tidy

go build -o Acrawler ./main.go

if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

# 设置日志路径
LOG_FILE="./Acrawler.log"
echo "日志：$LOG_FILE"

chmod +x Acrawler
nohup ./Acrawler >>"$LOG_FILE" 2>&1 &

# 进程验证
echo "服务已启动，进程信息："
# shellcheck disable=SC2009
ps aux | grep -v grep | grep Acrawler
