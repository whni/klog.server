#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import hashlib
import time
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)

# url + method
api_url = "{}/api/0/config/course_comment".format(host)
method = HTTPMethod.POST
course_comment_params = [
    {
        "pid": "102030405060708090000001",
        "course_record_pid": "102030405060708090000001",
        "comment_person_type": "teacher",
        "comment_person_pid": "102030405060708090000001",
        "comment_ts": int(time.time()),
        "comment_body": "comment for record 102030405060708090000001 teacher 1"
    },
    {
        "pid": "102030405060708090000002",
        "course_record_pid": "102030405060708090000001",
        "comment_person_type": "teacher",
        "comment_person_pid": "102030405060708090000002",
        "comment_ts": int(time.time()),
        "comment_body": "comment for record 102030405060708090000001 teacher 2"
    },
    {
        "pid": "102030405060708090000003",
        "course_record_pid": "102030405060708090000002",
        "comment_person_type": "teacher",
        "comment_person_pid": "102030405060708090000001",
        "comment_ts": int(time.time()),
        "comment_body": "comment for record 102030405060708090000002 teacher 1"
    },
    {
        "pid": "102030405060708090000004",
        "course_record_pid": "102030405060708090000002",
        "comment_person_type": "relative",
        "comment_person_pid": "102030405060708090000001",
        "comment_ts": int(time.time()),
        "comment_body": "comment for record 102030405060708090000001 relative 1"
    }
]

for params in course_comment_params:
    print("[create course comment: course record pid ({})]".format(params["course_record_pid"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()