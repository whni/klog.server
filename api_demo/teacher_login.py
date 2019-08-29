#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import hashlib
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)
        
# url + method
api_url = "{}/api/0/workflow/teacher/login".format(host)
method = HTTPMethod.POST
params = {
    "teacher_uid": "uid-usa-1001",
    "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()