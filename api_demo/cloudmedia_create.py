#!/usr/bin/python3

from http_request import HTTPRequest
from http_request import HTTPMethod
import sys
import os
import json
import random
from host_url import host_url_maker
        
# get host url
host = host_url_maker(sys.argv)

# url + method
api_url = "{}/api/0/config/cloudmedia".format(host)
method = HTTPMethod.POST

media_names = ["student1_image1.gif", "student1_image2.gif", "student1_image3.gif", "student1_image4.gif", "student1_image5.gif"]
#media_names = ["student2_video1.mp4", "student2_video2.mp4", "student2_video3.mp4", "student2_video4.mp4"]
for media_name in media_names:
    params = {
        "student_pid": "102030405060708090000001",
        "course_record_pid": "102030405060708090000001",
        "media_type": "image",
        "media_name": media_name,
        "media_url": "https://klogresourcediag159.blob.core.windows.net/klog-cloud-media/{}".format(media_name),
        "rank_score": random.uniform(50, 100),
        "media_tags": ["sport", "running"]
    }

    print("[create cloudmedia for {}]".format(media_name))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()
