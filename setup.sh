#!/bin/bash
. .env

PRIVATE_IP="127.0.0.1"

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/sendgrid -d $SENDGRID_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/mandrill -d $MANDRILL_API_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm -d $FCM_SERVER_KEY
curl -X PUT --url $PRIVATE_IP:8500/v1/kv/jwt_secret -d $JWT_SECRET
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

curl -X PUT --url $PRIVATE_IP:8500/v1/kv/fcm_credentials \
    -d '{
        "type": "service_account",
        "project_id": "yggdrasil-216009",
        "private_key_id": "a653a2d1522e10d34401cee6e46b217fad6b0db6",
        "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDvGHm6DNJ9ibpL\nNjDwA+QvEf4WQsgVwlj6CCFb2jsaW6y0KpzDSGKZRNOgEnQQq8/q1+LjgTzZK06V\nZ+hFONYNutulDXiY3gTwiRD57DTvnNeLx7rhxUnmhE0B6ru1hrnLRlGbhux7N20L\n7/O6fsduhW2sg2Xabss1CvfwnEN0TE82hUgBmWa8CPpHM5/A424W/mh+2q/DjB+F\n5XwihQy74qkb03ExbCjjcRH0r/lMvJD2uiJ/BlaFUPyfOY9EXxQ6AbdgP1I+V5Lw\nFzR1T9opbThfLgxE7JpeVp83sIteH0DHQPtAWWcjS9pjEaO0Jtat7fu5oka4b7y0\nNLScxYe1AgMBAAECggEADqHaEz2FJTsoTEGLAalPTeUlPmIEYjaEYPrd5cPngY+y\nNE12TgowRI5+gAmZEksdfJMNNtyaJ3U7yC91esUFbo4ssn9uTbaqqTjOMelXfQGF\nfG+y+22qPeIDX3Zo2e1ekfbh+Clr8Ad+lDAxY4yuMlCWm2voMBO/OmYLMe2yQBvs\nI0yaanFElRvrOs9YpbHOZJe4UsATQ2lHWFsdmnxYqh2MjxTDu7qKlbfRpGU/1G49\nDSZAQQlVEPNAQuuwATK8Zk9Kv+gIytNjWyMpp74V659R4XqGOo1KG6+ltvEy8B8Z\n8tXjsu/i6Z/Ma39pr5lr2NgTGxwSN8R0022uojoBIQKBgQD9+9h1CaO5W1krJJus\nx6pGK7AlD3q9wTm3qKX/kGs+RkloRqDcntSMvPh04gYdnU1BBRBu1gbXUKj3zwYC\nILvX8xNO31fr52A2HXY/94Klsbhz5QfcFw+4bKd3UwaxUswcsfzpGWWVqeCPaby8\nUu/kJby6mSL9Tbnsdlsv93WsVQKBgQDw/l+sdg6ZxjSxvlSEX0Wg+e5h6kBFBNMR\nSc6HXgEIDRDfY+3llltFpWqISHY9C0uqTOLqCiOycSEm1r8+ozMM+i5ojgClrrPy\nuOeYL4AnUlXOdpSxpntDQIs+IbiOTqlJ+CN348ynNAYreL5yy9lk7XYUpIfyUTcw\nEklHHvXN4QKBgCpdvDmppfnhVvpvXNhxZeHWn8dO0badaLdOFoKO4JS+vLf8MBEd\nHW9shjVZDpQaDSzvX1JduT5pYgYULnhkZXEcRvg9ITlkmxPrzjHVY6GoB5Ctf6Yi\n4qhs13J8Ev25RfrzZbCsX9mbQK4rgSQY4ZM5CEZIDSIbuQvMomFZ8fMRAoGADbof\njv5GhKHSnJE/1S5sE+hImNE0CfplX2W52uIS4veDj4XsphgzaLssX0jpWz7Kd4/b\nmQMs11w0TDcNO68mGEYj4Ja+GLHj1B6OKpACF7tL4e/FNn1KJiGGDOr0zT5TzD/T\nHWAfZyLCezVse7N61ZHlGPXqPHY904InJGLyE2ECgYEAx6UEOy+rcyDnoUhiZhjb\nQrdWIOzgDigO6MgFS1GI7POqEqpri+AF1kKCSNuMc7Ve7+f9lkDl+8+2OziIpjuO\ne91z2Qo01hTUyA7cFlSzQag7bmeOvGVCSmtird8G/HFrOSEg3YromSoIEUngNEDU\nEBfC26wknDQnxRvjSk4TdwU=\n-----END PRIVATE KEY-----\n",
        "client_email": "firebase-adminsdk-cdimi@yggdrasil-216009.iam.gserviceaccount.com",
        "client_id": "117261465157133358282",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-cdimi%40yggdrasil-216009.iam.gserviceaccount.com"
    }'

