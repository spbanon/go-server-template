#!/bin/sh

( rabbitmqctl wait --timeout 60 $RABBITMQ_PID_FILE ; \
rabbitmqctl add_user $RABBITMQ_USER $RABBITMQ_PASSWORD 2>/dev/null ; \
rabbitmqctl set_user_tags $RABBITMQ_USER administrator ; \
rabbitmqctl set_permissions -p / $RABBITMQ_USER  ".*" ".*" ".*" ; \
rabbitmqctl add_user $RABBITMQ_LOG_USER $RABBITMQ_LOG_PASSWORD 2>/dev/null ; \
rabbitmqctl set_user_tags $RABBITMQ_LOG_USER monitoring ; \
rabbitmqctl set_permissions -p / $RABBITMQ_LOG_USER  ".*" ".*" ".*" ; \

echo "*** User '$RABBITMQ_USER' with password '$RABBITMQ_PASSWORD' completed. ***" ; \
echo "*** User '$RABBITMQ_LOG_USER' with password '$RABBITMQ_LOG_PASSWORD' completed. ***" ; \
echo "*** Log in the WebUI at port 15672 (example: http:/localhost:15672) ***") &

rabbitmq-server $@