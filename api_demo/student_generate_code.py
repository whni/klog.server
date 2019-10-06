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
    "pid": "102030405060708090000001",          # student pid
    "teacher_pid": "102030405060708090000001"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()
