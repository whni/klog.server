#!/usr/bin/python3

from enum import Enum
import requests 
import json

class HTTPMethod(Enum):
    GET = 1
    POST = 2
    PUT = 3
    DELETE = 4

HTTPMethodStringMap = {
    HTTPMethod.GET: "GET",
    HTTPMethod.POST: "POST",
    HTTPMethod.PUT: "PUT",
    HTTPMethod.DELETE: "DELETE",
}

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
        self.resp = None
        
    def send(self):
        req_handler = HTTPMethodMap[self.method]
        self.resp = None
        print("Send <Request>")
        print("HTTP Method:", HTTPMethodStringMap[self.method])
        if self.method == HTTPMethod.GET or self.method == HTTPMethod.DELETE:
            self.resp = req_handler(url=self.api_url, params=self.params, timeout=3)
            api_url_full = self.api_url
            delimiter = "?"
            for key, val in self.params.items():
                api_url_full += delimiter + key + "=" + str(val)
                delimiter = "&"
            print("URL:", api_url_full)
        elif self.method == HTTPMethod.POST or self.method == HTTPMethod.PUT:
            print("URL:", self.api_url)
            print("Params:")
            print(json.dumps(self.params, indent=4, ensure_ascii=False))
            self.resp = req_handler(url=self.api_url, json=self.params, timeout=3)
        else:
            print("Unsupported HTTP method: {}".format(self.method))
  
    def print_resp(self):
        if self.resp is None:
            print("No response received.")
        else:
            print("Received {}".format(self.resp))
            print(json.dumps(self.resp.json(), indent=4, ensure_ascii=False))