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
api_url = "{}/api/0/config/student".format(host)
method = HTTPMethod.POST

student_params = [
    {
        "pid": "102030405060708090000001",
        "student_name": "Thomas Hu",
        "student_image_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1.jpg",
        "parent_wxid": "",
        "parent_name": "",
        "phone_number": "",
        "email": "",
        "binding_code": "",
        "binding_expire": 0,
        "teacher_pid": "102030405060708090000001"
    },
    {
        "pid": "102030405060708090000002",
        "student_name": "Bruce Wang",
        "student_image_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student2.jpg",
        "parent_wxid": "wxid-test2",
        "parent_name": "Madison Beer",
        "phone_number": "777-888-9999",
        "email": "beer@google.com",
        "binding_code": "",
        "binding_expire": 0,
        "teacher_pid": "102030405060708090000002"
    },
    {
        "pid": "102030405060708090000003",
        "student_name": "Tiffiny Shawn",
        "student_image_url": "",
        "parent_wxid": "wxid-test3",
        "parent_name": "Skylar Grey",
        "phone_number": "777-888-9999",
        "email": "skylar@facebook.com",
        "binding_code": "",
        "binding_expire": 0,
        "teacher_pid": "102030405060708090000003"
    },
    {
        "pid": "102030405060708090000004",
        "student_name": "Gintama Y.",
        "student_image_url": "",
        "parent_wxid": "wxid-test4",
        "parent_name": "Autumn Mendes",
        "phone_number": "777-888-9999",
        "email": "autumn@xxx.com",
        "binding_code": "",
        "binding_expire": 0,
        "teacher_pid": "102030405060708090000004"
    }
]

for params in student_params:
    print("[create student {}]".format(params["student_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()