FROM ubuntu:14.04

# ��װ���̿�ִ���ļ����� myapp.go ���룩
COPY nice /app/
COPY js /app/js
COPY pages /app/pages
COPY fonts /app/fonts
COPY css /app/css
COPY entrypoint.sh /app/

# �����Զ��������
CMD /app/entrypoint.sh

# The entrypoint of lightvm will start everything
# under `/etc/service` as daemon