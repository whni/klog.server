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
""" type ChatRoom struct {
	PID           primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	PrimerUserPID primitive.ObjectID `json:"primer_pid" bson:"primer_pid"`
	SecondUserPID primitive.ObjectID `json:"second_pid" bson:"second_pid"`
	RoomTag      string           `json:"room_tag" bson:"room_tag"`
	RecordTS      int64              `json:"record_ts" bson:"record_ts"`
} """      
# url + method
api_url = "{}/api/0/config/chatroom".format(host)
method = HTTPMethod.POST

user_params = [
    {
        "pid": "102030405060708090000001",
        "primer_pid": "102030405060708090000001",
        "updated_ts": 1576387331,
        "room_tag": "m20e"
    },
    {
        "pid": "102030405060708090000002",
        "primer_pid": "102030405060708090000002",
        "updated_ts": 1576387331,
        "room_tag": "f30e"
    }
]

for params in user_params:
    print("[create chatroom {}]".format(params["room_tag"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()