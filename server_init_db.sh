#!/bin/sh

mongo --eval 'load("db_mgmt/create_db.js")'
mongo --eval 'load("db_mgmt/create_user.js")'

