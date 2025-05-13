# استفاده از تصویر پایه nginx
FROM nginx:latest

# به‌روزرسانی بسته‌ها و نصب iproute2 (برای دستور ip addr)
RUN apt update && apt install -y iproute2 && apt clean

# کپی کردن اسکریپت اجرای اولیه
COPY start.sh /start.sh
RUN chmod +x /start.sh

# اجرای اسکریپت به‌عنوان entrypoint
ENTRYPOINT ["/start.sh"]
