#!/bin/bash
. .env

PRIVATE_IP="127.0.0.1"

docker exec -t postgres sh -c "psql -U postgres -c 'CREATE DATABASE scheduler'"

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sendgrid -d $SENDGRID_API_KEY 
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/mandrill -d $MANDRILL_API_KEY 
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm -d $FCM_SERVER_KEY 
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/twilio \
    -d '{
    "sid":"'$TWILIO_SID'",
    "key":"'$TWILIO_KEY'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/stripe \
    -d '{
    "public":"'$STRIPE_API_PUBLIC'",
    "secret":"'$STRIPE_API_SECRET'"
}'
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sql \
    -d '{
    "user": "admin",
    "password": "123"
}'