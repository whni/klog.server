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
api_url = "{}/api/0/config/teacher".format(host)
method = HTTPMethod.POST

teacher_params = [
    {
        "pid": "102030405060708090000001",
        "teacher_uid": "uid-usa-1001",
        "teacher_name": "Nicole Taylor",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "phone_number": "123-456-9876",
        "email": "nigoo@klog.com",
        "institute_pid": "102030405060708090000001"
    },
    {
        "pid": "102030405060708090000002",
        "teacher_uid": "uid-usa-1002",
        "teacher_name": "Wayne Grace",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "phone_number": "123-456-9876",
        "email": "wayne@klog.com",
        "institute_pid": "102030405060708090000001"
    },
    {
        "pid": "102030405060708090000003",
        "teacher_uid": "uid-usa-1003",
        "teacher_name": "Fantasy God",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
        "phone_number": "000-111-2222",
        "email": "fanfan@klog.com",
        "institute_pid": "102030405060708090000002"
    },
    {
        "pid": "102030405060708090000004",
        "teacher_uid": "uid-usa-1004",
        "teacher_name": "倪炜恒",
        "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
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

# Query teacher
api_url = "{}/api/0/config/teacher?pid=all&fkey=institute_pid&fid=102030405060708090000002".format(host)
method = HTTPMethod.GET
params = {}

print("[Query teacher]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()