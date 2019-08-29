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
api_url = "{}/api/0/config/institute".format(host)
method = HTTPMethod.POST

institute_params = [
    {
        "pid": "102030405060708090000001",
        "institute_uid": "uid-usa-0001",
        "institute_name": "Institute 1",
        "address": {
            "street": "180 Elm Ct",
            "code": "94086",
            "city": "Sunnyvale",
            "state": "CA",
            "country": "USA"
        }
    },
    {
        "pid": "102030405060708090000002",
        "institute_uid": "uid-usa-0002",
        "institute_name": "Institute 2",
        "address": {
            "street": "Valley Green 6",
            "code": "95014",
            "city": "Cupertino",
            "state": "CA",
            "country": "USA"
        }
    }
]


for params in institute_params:
    print("[create institute {}]".format(params["institute_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()