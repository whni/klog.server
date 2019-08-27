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

# get student info
api_url = "{}/api/0/workflow/parent/findstudent".format(host)
method = HTTPMethod.POST
params = {
    "parent_wxid": "wxid-0123456789",
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

student = http_req.resp.json()["payload"]

# get cloud media for student
api_url = "{}/api/0/workflow/student/mediaquery".format(host)
method = HTTPMethod.POST
params = {
    "student_pid": student["pid"],
    "start_ts": 0,
    "end_ts": 1566945630,
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()