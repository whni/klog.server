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
api_url = "{}/api/0/config/relative".format(host)
method = HTTPMethod.POST

relative_params = [
    {
        "pid": "102030405060708090000001",
        "relative_name": "Relative 1",
        "relative_wxid": "relative_wxid_1",
        "phone_number": "123-456-7890",
        "email": "relative_1@gmail.com"
    },
    {
        "pid": "102030405060708090000002",
        "relative_name": "Relative 2",
        "relative_wxid": "relative_wxid_2",
        "phone_number": "987-654-1230",
        "email": "relative_2@gmail.com"
    }
]


for params in relative_params:
    print("[create relative {}]".format(params["relative_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()