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
api_url = "{}/api/0/workflow/relative/findstudent".format(host)
method = HTTPMethod.POST
params = {
    "relative_wxid": "relative_wxid_2",
}
print("[find student for relative]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

references = http_req.resp.json()["payload"]
reference = None
if len(references) > 0:
    reference = references[0]
else:
    print("No student-relative references found for wechat id: {}".format(params["relative_wxid"]))
    exit(0)

# edit student-relative info as student-main-relative
api_url = "{}/api/0/workflow/relative/extra/edit".format(host)
method = HTTPMethod.POST
params = {
    "student_pid":          "102030405060708090000002",
    "relative_wxid":        "relative_wxid_2",
    "sec_relative_wxid":    "relative_wxid_3",
    "relationship":         "uncle",
    "operation":            "add",
}
print("[add student and relative record]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

api_url = "{}/api/0/workflow/relative/extra/edit".format(host)
method = HTTPMethod.POST
params = {
    "student_pid":          "102030405060708090000002",
    "relative_wxid":        "relative_wxid_2",
    "sec_relative_wxid":    "relative_wxid_3",
    "relationship":         "uncle",
    "operation":            "delete",
}
print("[add student and relative record]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

# get cloud media for student
api_url = "{}/api/0/workflow/student/mediaquery".format(host)
method = HTTPMethod.POST
params = {
    "student_pid": reference["student_pid"],
    "start_ts": 0,
    "end_ts": int(time.time()),
}

print("[find cloud media for student]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()

# get story for student
api_url = "{}/api/0/workflow/student/storyquery".format(host)
method = HTTPMethod.POST
params = {
    "student_pid": reference["student_pid"],
    "start_ts": 0,
    "end_ts": int(time.time()),
}

print("[find story for student]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()