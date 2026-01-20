# Setup RabbitMQ

## Setup SOM User

Run as root:

```
rabbitmqctl list_users
echo "somMsgBusPassword" | rabbitmqctl add_user som

rabbitmqctl set_permissions som "som..*" "som..*" "som..*"
rabbitmqctl set_user_tags som "monitoring"
```

```
echo "somMsgBusPassword" | rabbitmqctl change_password som
```

## Commands

### User

```
rabbitmqctl list_users
rabbitmqctl add_user vogtp

rabbitmqctl set_permissions vogtp ".*" ".*" ".*"
rabbitmqctl set_user_tags vogtp admin
```
