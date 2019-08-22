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
api_url = "http://{}:8080/api/0/workflow/student/binding".format(host)
method = HTTPMethod.POST
params = {
    "parent_wxid": "wxid-my-test",
    "parent_name": "For Test",
    "phone_number": "777-888-9999",
    "email": "test@klog.com",
    "binding_code": "blfe488fggckigkkbio0",
}

http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()