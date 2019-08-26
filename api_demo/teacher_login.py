#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
        
# url + method
host = "127.0.0.1:80"
if len(sys.argv) > 1:
    host = sys.argv[1]
api_url = "http://{}/api/0/workflow/teacher/login".format(host)
method = HTTPMethod.POST
params = {
    "teacher_uid": "uid-usa-1001",
    "teacher_key": "no_key"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()