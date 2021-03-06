FROM scratch

ENV MONEYIO_LOCAL_HOST 0.0.0.0
ENV MONEYIO_LOCAL_PORT 1666
ENV MONEYIO_LOG_LEVEL 0

ENV MYSQL_SERVICE_HOST localhost
ENV MYSQL_SERVICE_PORT 3306
ENV MYSQL_SERVICE_DATABASE moneyio
ENV MYSQL_SERVICE_USERNAME moneyio
ENV MYSQL_SERVICE_PASSWORD password

ENV MONEYIO_TELEGRAM_TOKEN 123454352:tokenhere

EXPOSE $MONEYIO_LOCAL_PORT

COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/money.io /

CMD ["/money.io"]
