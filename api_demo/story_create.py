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
api_url = "{}/api/0/config/story".format(host)
method = HTTPMethod.POST
story_params = [
    {
        "student_pid": "102030405060708090000001",
        "story_ts": 1576387331,
        "story_template": {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_content": "template_food.zip",
                "clip_duration": 5,
                "clip_sequence": 1,
                "type": "json"
            },
            {
                "clip_content": "user_url1",
                "clip_duration": 4,
                "clip_sequence": 2,
                "type": "video"
            },
            {
                "clip_content": "user_url2",
                "clip_duration": 4,
                "clip_sequence": 3,
                "type": "video"
            },
            {
                "clip_content": "user_url3",
                "clip_duration": 5,
                "clip_sequence": 4,
                "type": "video"
            }
        ],
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food.zip",
        "template_mp4_movie": "nul",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food_1.mp3",
        "template_name": "001_Food_Delightful_short",
        "template_tags": [
            "happy",
            "kids",
            "family",
            "food"
        ]
    }
    },
    {
        "student_pid": "102030405060708090000002",
        "story_ts": 1576387531,
        "story_template": {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 6,
        "template_clip_time_content": [
            {
                "clip_content": "template_food.zip",
                "clip_duration": 5,
                "clip_sequence": 1,
                "type": "json"
            },
            {
                "clip_content": "user_url1",
                "clip_duration": 4,
                "clip_sequence": 2,
                "type": "video"
            },
            {
                "clip_content": "user_url2",
                "clip_duration": 4,
                "clip_sequence": 3,
                "type": "video"
            },
            {
                "clip_content": "user_url3",
                "clip_duration": 5,
                "clip_sequence": 4,
                "type": "video"
            },
            {
                "clip_content": "user_url4",
                "clip_duration": 5,
                "clip_sequence": 5,
                "type": "video"
            },
            {
                "clip_content": "user_url5",
                "clip_duration": 4,
                "clip_sequence": 6,
                "type": "video"
            }
        ],
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food.zip",
        "template_mp4_movie": "",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food_3.mp3",
        "template_name": "001_Food_Delightful_long",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    }
    }
]

for params in story_params:
    print("[create store for student {}]".format(params["student_pid"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()

# Query student
api_url = "{}/api/0/config/story?pid=all".format(host)
method = HTTPMethod.GET
params = {}

print("[Query story]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()