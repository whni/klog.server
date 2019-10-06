#!/usr/bin/python3
        
def host_url_maker(sys_arg):
    host_map = {
        "local": "http://127.0.0.1:80",
        "remote": "https://klogservice.westus2.cloudapp.azure.com:443",
    }
    if len(sys_arg) == 1:
        host_type = "local"
    if len(sys_arg) > 1:
        host_type = sys_arg[1]
        if host_type != "remote":
            host_type = "local"
    return host_map[host_type]
