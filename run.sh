#!/bin/bash

go mod tidy

go build -o Acrawler cmd/crawler/main.go

if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

# 设置日志路径（建议使用用户有写入权限的路径）
LOG_FILE="./Acrawler.log"  # 可修改为 /var/log/Acrawler.log（需确保有权限）
echo "日志将输出到：$LOG_FILE"

# 添加执行权限
chmod +x Acrawler

# 后台运行阶段
echo "启动后台进程..."
nohup ./Acrawler > "$LOG_FILE" 2>&1 &

# 进程验证
echo "服务已启动，进程信息："
ps aux | grep -v grep | grep Acrawler

# 使用提示
echo -e "\n操作完成！您可以通过以下命令查看日志："
echo "tail -f $LOG_FILE"