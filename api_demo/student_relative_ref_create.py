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
api_url = "{}/api/0/config/reference/student_relative".format(host)
method = HTTPMethod.POST
reference_params = [
    {
        "pid": "102030405060708090000001",
        "student_pid": "102030405060708090000001",
        "relative_pid": "102030405060708090000001",
        "relationship": "mother",
        "is_main": True
    },
    {
        "pid": "102030405060708090000002",
        "student_pid": "102030405060708090000001",
        "relative_pid": "102030405060708090000002",
        "relationship": "father",
        "is_main": False
    },
    {
        "pid": "102030405060708090000003",
        "student_pid": "102030405060708090000001",
        "relative_pid": "102030405060708090000003",
        "relationship": "grandpa",
        "is_main": False
    },
    {
        "pid": "102030405060708090000004",
        "student_pid": "102030405060708090000002",
        "relative_pid": "102030405060708090000005",
        "relationship": "falther",
        "is_main": True
    },
    {
        "pid": "102030405060708090000005",
        "student_pid": "102030405060708090000002",
        "relative_pid": "102030405060708090000004",
        "relationship": "mother",
        "is_main": False
    }
]

for params in reference_params:
    print("[create student ({}) - relative ({}) reference]".format(params["student_pid"], params["relative_pid"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()