#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import hashlib
        
# url + method
host = "http://127.0.0.1:80"
if len(sys.argv) > 1:
    host = sys.argv[1]
api_url = "{}/api/0/config/teacher".format(host)
method = HTTPMethod.POST
params = {
    "teacher_uid": "uid-for-test",
    "teacher_key": hashlib.sha256("test_password".encode()).hexdigest(),
    "teacher_name": "Test Teacher",
    "class_name": "Summer",
    "phone_number": "123-456-9876",
    "email": "tttt@klog.com",
    "institute_pid": "102030405060708090000001"
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()