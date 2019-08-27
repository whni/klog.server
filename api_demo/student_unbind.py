#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
        
# url + method
host = "http://127.0.0.1:80"
if len(sys.argv) > 1:
    host = sys.argv[1]
api_url = "{}/api/0/workflow/student/unbind".format(host)
method = HTTPMethod.POST
params = {
    "pid": "102030405060708090000001",
    "parent_wxid": "wxid-0123456789",
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()