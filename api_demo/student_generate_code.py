#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
        
# url + method
host = "127.0.0.1"
if len(sys.argv) > 1:
    host = sys.argv[1]
api_url = "http://{}:8080/api/0/workflow/student/generatecode".format(host)
method = HTTPMethod.POST
params = {
    "pid": "102030405060708090000001",
    "teacher_pid": "102030405060708090000001"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()
