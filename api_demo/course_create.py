#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import hashlib
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)
        
# url + method
api_url = "{}/api/0/config/course".format(host)
method = HTTPMethod.POST

course_params = [
    {
        "pid": "102030405060708090000001",
        "course_uid": "uid-course-1001",
        "course_name": "Advanced Math",
        "course_intro": "advanced math course",
        "course_targets": [
            {
                "tag": "c1",
                "desc": "1st course target"
            },
            {
                "tag": "c2",
                "desc": "2nd course target"
            }
        ],
        "institute_pid": "102030405060708090000001",
        "teacher_pid": "102030405060708090000001",
        "assistant_pid": "102030405060708090000002"
    },
    {
        "pid": "102030405060708090000002",
        "course_uid": "uid-course-1002",
        "course_name": "English",
        "course_intro": "New English",
        "course_targets": [
            {
                "tag": "c1",
                "desc": "1st course target"
            },
            {
                "tag": "c2",
                "desc": "2nd course target"
            }
        ],
        "institute_pid": "102030405060708090000001",
        "teacher_pid": "102030405060708090000002",
        "assistant_pid": "102030405060708090000003"
    }
]

for params in course_params:
    print("[create course {}]".format(params["course_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()