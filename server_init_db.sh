#!/bin/sh

mongo --eval 'load("db_mgmt/create_db.js")'
mongo --eval 'load("db_mgmt/create_user.js")'

python3 api_demo/institute_create.py
python3 api_demo/teacher_create.py
python3 api_demo/student_create.py
python3 api_demo/cloudmedia_create.py

python3 api_demo/student_bind.py
python3 api_demo/parent_student_query.py

