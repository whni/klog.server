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
api_url = "{}/api/0/config/user".format(host)
method = HTTPMethod.POST

user_params = [
    {
        "pid": "102030405060708090000001",
        "user_email": "test1@gmail.com",
        "updated_ts": 1576387331,
        "user_description": "source from fb"
    },
    {
        "pid": "102030405060708090000002",
        "user_email": "test2@gmail.com",
        "updated_ts": 1576388331,
        "user_description": "source from fb"
    }
]

for params in user_params:
    print("[create user {}]".format(params["user_email"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()