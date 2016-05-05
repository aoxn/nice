FROM 61.160.36.122:8080/lightvm:latest

# 安装进程可执行文件（由 myapp.go 编译）
COPY nice /app/
COPY js /app/js
COPY pages /app/pages
COPY fonts /app/fonts
COPY css /app/css
COPY entrypoint.sh /app/

# 设置自动拉起进程
RUN mkdir /etc/service/nice
COPY entrypoint.sh /etc/service/nice/run
RUN chmod +x /etc/service/nice/run

# The entrypoint of lightvm will start everything
# under `/etc/service` as daemon