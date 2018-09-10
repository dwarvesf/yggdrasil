#!/bin/bash
. .env

sleep 2

PRIVATE_IP="127.0.0.1"

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sendgrid -d $SENDGRID_API_KEY 
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/twilio \
    -d '{
    "sid":"'$TWILIO_SID'",
    "key":"'$TWILIO_KEY'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sql \
    -d '{
    "user": "admin",
    "password": "123"
}'