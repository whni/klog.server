#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)

# generate code
api_url = "{}/api/0/workflow/student/generatecode".format(host)
method = HTTPMethod.POST
params = {
    "student_pid": "102030405060708090000002",          # student pid
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

# binding
binding_code = http_req.resp.json()["payload"]["binding_code"]
api_url = "{}/api/0/workflow/student/bind".format(host)
method = HTTPMethod.POST
params = {
    "relative_wxid": "relative_wxid_2",
    "relationship": "Mother",
    "binding_code": binding_code
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()
