#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)

# url + method
api_url = "{}/api/0/workflow/parent/wxlogin".format(host)
method = HTTPMethod.POST
params = {
    "appid": "wx03932c08a933f9a9",
    "js_code": "001vr5Ac15ZYGw0D5mwc1j4Wzc1vr5Af",
    "secret": "70a5b49cfcf97b31dbfc9a136e2295b1"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()