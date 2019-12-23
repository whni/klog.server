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
api_url = "{}/api/0/config/reference/student_course".format(host)
method = HTTPMethod.POST
reference_params = [
    {
        "pid": "102030405060708090000001",
        "student_pid": "102030405060708090000001",
        "course_pid": "102030405060708090000001",
    },
    {
        "pid": "102030405060708090000002",
        "student_pid": "102030405060708090000001",
        "course_pid": "102030405060708090000002",
    },
    {
        "pid": "102030405060708090000003",
        "student_pid": "102030405060708090000002",
        "course_pid": "102030405060708090000001",
    },
    {
        "pid": "102030405060708090000004",
        "student_pid": "102030405060708090000002",
        "course_pid": "102030405060708090000002",
    },
    {
        "pid": "102030405060708090000005",
        "student_pid": "102030405060708090000003",
        "course_pid": "102030405060708090000002",
    },
    {
        "pid": "102030405060708090000006",
        "student_pid": "102030405060708090000004",
        "course_pid": "102030405060708090000002",
    }
]

for params in reference_params:
    print("[create student ({}) - course ({}) reference]".format(params["student_pid"], params["course_pid"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()