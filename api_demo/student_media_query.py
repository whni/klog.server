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

# get all student info
api_url = "{}/api/0/config/student?pid=all".format(host)
method = HTTPMethod.GET
params = {}

print("[Query all students]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
#http_req.print_resp()
students = http_req.resp.json()["payload"]

for student in students :
    # do you need to add filter here ? filter videos without any tags
    api_url = "{}/api/0/workflow/student/mediaquery".format(host)
    method = HTTPMethod.POST
    # you can change the query time range
    params = {
        "student_pid": student["pid"],
        "start_ts": 0,
        "end_ts": int(time.time()),
    }
    print("find cloud media for student " + student["pid"])
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    cloudMedias = http_req.resp.json()["payload"]
    for cloudMedia in cloudMedias:
        # check tag and type
        if cloudMedia["media_type"] == "video" :
            # download media, update tags and rankscore here
            print(cloudMedia["media_url"])
            # cloudMedia["rank_score"] = score1
            # cloudMedia["media_tags"].append("tag1")
            # Update cloudMedia into db
            api_url = "{}/api/0/config/cloudmedia".format(host)
            method = HTTPMethod.PUT
            params = cloudMedia
            http_req = HTTPRequest(api_url, method, params)
            http_req.send()
            http_req.print_resp()

