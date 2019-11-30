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
api_url = "{}/api/0/config/relative".format(host)
method = HTTPMethod.POST

relative_params = []
for i in range(1, 7):
    relative = {}
    relative["pid"] = "1020304050607080900000{:02d}".format(i)
    relative["relative_name"] = "Relative {}".format(i)
    relative["relative_wxid"] = "relative_wxid_{}".format(i)
    relative["phone_number"] = "123-123-12{:02d}".format(i)
    relative["email"] = "relative_{}@klog.com".format(i)
    relative_params.append(relative)


for params in relative_params:
    print("[create relative {}]".format(params["relative_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()
