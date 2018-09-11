#!/bin/bash
. .env

PRIVATE_IP="127.0.0.1"

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sendgrid -d $SENDGRID_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/mandrill -d $MANDRILL_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm -d $FCM_SERVER_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/project_id -d $PROJECT_ID
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