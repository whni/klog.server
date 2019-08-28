#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import time
        
# url + method
host = "http://127.0.0.1:80"
if len(sys.argv) > 1:
    host = sys.argv[1]

# generate code
api_url = "{}/api/0/workflow/student/mediaquery".format(host)
method = HTTPMethod.POST
params = {
    "student_pid": "1020304050607080900000ff",      # find media for this student pid
    "start_ts": 0,                                  # start unix timestamp (from 0)
    "end_ts": int(time.time()),                     # end unix timestamp (to current time)
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()
