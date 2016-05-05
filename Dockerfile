FROM ubuntu:14.04

# 安装进程可执行文件（由 myapp.go 编译）
COPY nice /app/
COPY js /app/js
COPY pages /app/pages
COPY fonts /app/fonts
COPY css /app/css
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

# 设置自动拉起进程
CMD /app/entrypoint.sh

# The entrypoint of lightvm will start everything
# under `/etc/service` as daemon