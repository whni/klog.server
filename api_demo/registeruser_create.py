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
""" type RegisterUser struct {
	PID          primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	UserEmail    string             `json:"user_email" bson:"user_email"`
	UserPassWord string             `json:"user_password" bson:"user_password"`
	Dob          UserDOB            `json:"dob" bson:"dob"`
	Address      AddressInfo        `json:"address" bson:"address"`
	Gender       string             `json:"gender" bson:"gender"`
	UpdatedTS    int64              `json:"updated_ts" bson:"updated_ts"`
}   """      
# url + method
api_url = "{}/api/0/config/registeruser".format(host)
method = HTTPMethod.POST

user_params = [
    {
        "pid": "102030405060708090000001",
        "user_email": "test1@gmail.com",
        "updated_ts": 1576387331,
        "user_password": "axbcdxes",
        "gender":"male",
        "dob":{
            "year": 1980,
            "month":4,
            "day":10
        },
        "address": {
            "street": "180 Elm Ct",
            "code": "94086",
            "city": "Sunnyvale",
            "state": "CA",
            "country": "USA"
        }
    },
    {
        "pid": "102030405060708090000002",
        "user_email": "test2@gmail.com",
        "updated_ts": 1576388331,
        "user_password": "axbcdxes123",
        "address": {
            "street": "180 Elm Ct",
            "code": "94086",
            "city": "Sunnyvale",
            "state": "CA",
            "country": "USA"
        },
        "dob":{
            "year": 1980,
            "month":4,
            "day":10
        },
        "gender":"male"
    },
]

for params in user_params:
    print("[create register user {}]".format(params["user_email"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()

api_url = "{}/api/0/workflow/registeruser/login".format(host)
method = HTTPMethod.POST

for params in user_params:
    print("[ login register user {}]".format(params["user_email"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()