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

#media_names1 = ["student1_video1.mp4", "student1_video2.mp4", "student1_video3.mp4", "student1_video4.mp4", "student1_video5.mp4", "student1_video5.mp4"]
#media_names2 = ["student2_video1.mp4", "student2_video2.mp4", "student2_video3.mp4", "student2_video4.mp4"]
media_name_array = [10, 4, 10, 4]
student_pids = ["102030405060708090000001", "102030405060708090000002","102030405060708090000003","102030405060708090000004"]
course_record_pids = ["102030405060708090000001", "102030405060708090000003", "102030405060708090000005", "102030405060708090000006"]

for i in range(len(media_name_array)):
    for j in range(1, media_name_array[i]+1):
        media_name = "student{}_video{}.mp4".format(i+1,j)
        params = {
            "student_pid": student_pids[i],
            "course_record_pid": course_record_pids[i] if i == 0 else "000000000000000000000000",
            "media_type": "video",
            "media_name": media_name,
            "media_url": "https://klogresourcediag159.blob.core.windows.net/klog-cloud-media/{}".format(media_name),
            "media_tags": []
           # "media_tags": [{"tag_name":"test","tag_score":0.789}]
        }

        print("[create cloudmedia for {}]".format(media_name))
        http_req = HTTPRequest(api_url, method, params)
        http_req.send()
        http_req.print_resp()
