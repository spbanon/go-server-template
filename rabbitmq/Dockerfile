FROM rabbitmq:3.11.2-management

ENV RABBITMQ_USER=*****
ENV RABBITMQ_PASSWORD=*****

ENV RABBITMQ_PID_FILE=/var/lib/rabbitmq/mnesia/rabbitmq

ADD init.sh /init.sh
RUN chmod +x /init.sh
EXPOSE 15672

CMD ["/init.sh"]