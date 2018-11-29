#!/bin/bash
. .env

PRIVATE_IP="127.0.0.1"

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sendgrid -d $SENDGRID_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/mandrill -d $MANDRILL_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm -d $FCM_SERVER_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/jwt_secret -d $JWT_SECRET
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/mailgun \
    -d '{
    "api_key":"'$MAILGUN_API_KEY'",
    "pub_key":"'$MAILGUN_PUB_KEY'",
    "domain":"'$MAILGUN_DOMAIN'",
    "host":"'$MAILGUN_HOSTNAME'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/twilio \
    -d '{
    "sid":"'$TWILIO_SID'",
    "key":"'$TWILIO_KEY'",
    "token":"'$TWILIO_TOKEN'",
    "number":"'$TWILIO_NUMBER'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/stripe \
    -d '{
    "public":"'$STRIPE_API_PUBLIC'",
    "secret":"'$STRIPE_API_SECRET'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/db-scheduler \
    -d '{
    "user": "admin",
    "password": "123",
    "db": "scheduler"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/db-identity \
    -d '{
    "user": "admin",
    "password": "123",
    "db": "identity"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/db-organization \
    -d '{
    "user": "admin",
    "password": "123",
    "db": "organization"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/kafka \
    -d '{
    "address": "10.5.0.6",
    "port": 9092
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/db-networks \
    -d '{
    "user": "admin",
    "password": "123",
    "db": "networks"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm_credentials -d $FCM_CREDENTIAL
