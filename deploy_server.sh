#!/bin/sh

sudo killall klog.server

mongo --eval 'load("db_mgmt/create_db.js")'
mongo --eval 'load("db_mgmt/create_user.js")'

sudo ./klog.server &

sleep 3

python3 api_demo/cloudmedia_create.py

