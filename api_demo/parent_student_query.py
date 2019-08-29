#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import time
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)
        
# get student info
api_url = "{}/api/0/workflow/parent/findstudent".format(host)
method = HTTPMethod.POST
params = {
    "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
}
print("[find student for parent]")
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
    "end_ts": int(time.time()),
}

print("[find cloud media for student]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()