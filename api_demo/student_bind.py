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

# generate code
api_url = "{}/api/0/workflow/student/generatecode".format(host)
method = HTTPMethod.POST
params = {
    "pid": "102030405060708090000001",          # student pid
    "teacher_pid": "102030405060708090000001"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

# binding
binding_code = http_req.resp.json()["payload"]["binding_code"]
api_url = "{}/api/0/workflow/student/bind".format(host)
method = HTTPMethod.POST
params = {
    "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
    "parent_name": "Bruce Wayne",
    "phone_number": "777-888-9999",
    "email": "bruce@klog.com",
    "binding_code": binding_code
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()