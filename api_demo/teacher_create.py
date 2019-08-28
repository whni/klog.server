#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import hashlib
        
# url + method
host = "http://127.0.0.1:80"
if len(sys.argv) > 1:
    host = sys.argv[1]
api_url = "{}/api/0/config/teacher".format(host)
method = HTTPMethod.POST

teacher_params = [
    {
        "pid": "102030405060708090000001",
        "teacher_uid": "uid-usa-1001",
        "teacher_name": "Nicole Taylor",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "class_name": "GoldenEye",
        "phone_number": "123-456-9876",
        "email": "nigoo@klog.com",
        "institute_pid": "102030405060708090000001"
    },
    {
        "pid": "102030405060708090000002",
        "teacher_uid": "uid-usa-1002",
        "teacher_name": "Wayne Grace",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "class_name": "FastWind",
        "phone_number": "123-456-9876",
        "email": "wayne@klog.com",
        "institute_pid": "102030405060708090000001"
    },
    {
        "pid": "102030405060708090000003",
        "teacher_uid": "uid-usa-1003",
        "teacher_name": "Fantasy God",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "class_name": "CloudTop",
        "phone_number": "000-111-2222",
        "email": "fanfan@klog.com",
        "institute_pid": "102030405060708090000002"
    },
    {
        "pid": "102030405060708090000004",
        "teacher_uid": "uid-usa-1004",
        "teacher_name": "倪炜恒",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "class_name": "UnderWorld",
        "phone_number": "619-763-1020",
        "email": "summer@klog.com",
        "institute_pid": "102030405060708090000002"
    }
]

for params in teacher_params:
    print("[create teacher {}]".format(params["teacher_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()