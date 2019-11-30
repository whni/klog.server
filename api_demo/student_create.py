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
api_url = "{}/api/0/config/student".format(host)
method = HTTPMethod.POST
student_params = [
    {
        "pid": "102030405060708090000001",
        "student_name": "Thomas Hu",
        "student_image_name": "",
        "student_image_url": "",
        "binding_code": "",
        "binding_expire": 0
    },
    {
        "pid": "102030405060708090000002",
        "student_name": "Bruce Wang",
        "student_image_name": "",
        "student_image_url": "",
        "binding_code": "",
        "binding_expire": 0
    },
    {
        "pid": "102030405060708090000003",
        "student_name": "Tiffiny Shawn",
        "student_image_name": "",
        "student_image_url": "",
        "binding_code": "",
        "binding_expire": 0
    },
    {
        "pid": "102030405060708090000004",
        "student_name": "Gintama Y.",
        "student_image_name": "",
        "student_image_url": "",
        "binding_code": "",
        "binding_expire": 0
    }
]

for params in student_params:
    print("[create student {}]".format(params["student_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()