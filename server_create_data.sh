#!/bin/sh

python3 api_demo/institute_create.py
python3 api_demo/teacher_create.py
python3 api_demo/course_create.py
python3 api_demo/student_create.py
python3 api_demo/relative_create.py

python3 api_demo/student_relative_ref_create.py
python3 api_demo/student_course_ref_create.py

python3 api_demo/course_record_create.py
python3 api_demo/course_comment_create.py
python3 api_demo/cloudmedia_create.py

python3 api_demo/student_unbind.py
python3 api_demo/student_bind.py
python3 api_demo/relative_student_query.py

python3 api_demo/template_create.py
