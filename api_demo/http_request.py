#!/usr/bin/python3

from enum import Enum
import requests 
import json

class HTTPMethod(Enum):
    GET = 1
    POST = 2
    PUT = 3
    DELETE = 4

HTTPMethodMap = {
    HTTPMethod.GET: requests.get,
    HTTPMethod.POST: requests.post,
    HTTPMethod.PUT: requests.put,
    HTTPMethod.DELETE: requests.delete,
}
  
class HTTPRequest:
    def __init__(self, api_url, method, params):
        self.api_url = api_url
        self.method = method
        self.params = {key:val for key, val in params.items()}
        
    def send(self):
        req_handler = HTTPMethodMap[self.method]
        resp = None
        if self.method == HTTPMethod.GET or self.method == HTTPMethod.DELETE:
            resp = req_handler(url=self.api_url, params=self.params, timeout=3)
        elif self.method == HTTPMethod.POST or self.method == HTTPMethod.PUT:
            resp = req_handler(url=self.api_url, json=self.params, timeout=3)
        else:
            print("Unsupported HTTP method: {}".format(self.method))
        
        if resp is None:
            print("No response received.")
        else:
            print("Received {}".format(resp))
            print(json.dumps(resp.json(), indent=4, ensure_ascii=False))
  