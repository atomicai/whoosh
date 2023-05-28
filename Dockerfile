FROM rabbitmq:3-management AS rabbitmq
RUN rabbitmq-plugins enable rabbitmq_web_stomp