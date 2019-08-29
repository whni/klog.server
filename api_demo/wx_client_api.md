# WeChat Client API

For HTTPs API, the host DNS name is http://klogserver.westus2.cloudapp.azure.com:443.
Currently, all HTTP responses return results in JSON. HTTP status code 200 means a correct result, while other codes show an error appearing at server. HTTP response from server normally has "payload" and/or "message" fields: "payload" field in response JSON contains query results, and "message" field in response JSON addresses some additional information, mainly for explaining why error occurs.

To retrieve HTTP response results, please check HTTP status code first before reading any JSON field.

For API demo, please check some python script `api_demo` folder:
1. `parent_wxlogin.py`: parent wechat login
2. `student_bind.py`: student-parent bind
3. `student_unbind.py`: student-parent unbind
4. `parent_student_query.py`: parent query student information and student query cloudmedia


#### Parent WeChat login API
  - URL: /api/0/workflow/parent/wxlogin
  - Method: POST
  - Request JSON:

        {
            "appid": "wx03932c08a933f9a9",
            "js_code": "001vr5Ac15ZYGw0D5mwc1j4Wzc1vr5Af",
            "secret": "70a5b49cfcf97b31dbfc9a136e2295b1"
        }

  - Response JSON:

        {
            "payload": {
                "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ"
            }
        }
    or 

        {
            "message": "some information"
        }


#### Student-Parent Bind API
  - URL: /api/0/workflow/student/bind
  - Method: POST
  - Request JSON:

        {
            "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
            "parent_name": "Bruce Wayne",
            "phone_number": "777-888-9999",
            "email": "bruce@klog.com",
            "binding_code": "some binding_code related to one student"
        }
    
  - Response JSON (return student struct in payload):

        {
            "payload": {
                "pid": "102030405060708090000001",
                "student_name": "Thomas Hu",
                "student_image_name": "student1.jpg",
                "student_image_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1.jpg",
                "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
                "parent_name": "Bruce Wayne",
                "phone_number": "777-888-9999",
                "email": "bruce@klog.com",
                "binding_code": "",
                "binding_expire": 0,
                "teacher_pid": "102030405060708090000001"
            }
        }
    or 

        {
            "message": "some information"
        }


#### Student-Parent Unbind API
  - URL: /api/0/workflow/student/unbind
  - Method: POST
  - Request JSON (pid mean student pid):

        {
            "pid": "102030405060708090000001",
            "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
        }
    
  - Response JSON (return student struct in payload):

        {
            "payload": {
                "pid": "102030405060708090000001",
                "student_name": "Thomas Hu",
                "student_image_name": "student1.jpg",
                "student_image_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1.jpg",
                "parent_wxid": "",
                "parent_name": "",
                "phone_number": "",
                "email": "",
                "binding_code": "",
                "binding_expire": 0,
                "teacher_pid": "102030405060708090000001"
            }
        }
    or 

        {
            "message": "some information"
        }



#### Parent Query Student Information API
  - URL: /api/0/workflow/parent/findstudent
  - Method: POST
  - Request JSON:

        {
            "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ"
        }

  - Response JSON (return student struct in payload):

        {
            "payload": {
                "pid": "102030405060708090000001",
                "student_name": "Thomas Hu",
                "student_image_name": "student1.jpg",
                "student_image_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1.jpg",
                "parent_wxid": "orgQa44wYyOpdShmXAsHtSfjMjeQ",
                "parent_name": "Bruce Wayne",
                "phone_number": "777-888-9999",
                "email": "bruce@klog.com",
                "binding_code": "",
                "binding_expire": 0,
                "teacher_pid": "102030405060708090000001"
            }
        }
    or 

        {
            "message": "some information"
        }


#### Parent Query Student CloudMedia API
  - URL: /api/0/workflow/student/mediaquery
  - Method: POST
  - Request JSON (start_ts, end_ts: unix timestamp):

        {
            "student_pid": "102030405060708090000001",
            "start_ts": 0,
            "end_ts": 1567112011
        }

  - Response JSON (return student struct in payload):

        {
            "payload": [
                {
                    "pid": "5d67140850ea66aa1b0ffa75",
                    "media_type": "video",
                    "media_name": "student1_video1.mp4",
                    "media_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1_video1.mp4",
                    "rank_score": 87.08606384551757,
                    "create_ts": 1567015818,
                    "content_length": 8929558,
                    "student_pid": "102030405060708090000001"
                },
                {
                    "pid": "5d67140850ea66aa1b0ffa77",
                    "media_type": "video",
                    "media_name": "student1_video3.mp4",
                    "media_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1_video3.mp4",
                    "rank_score": 61.96497775137912,
                    "create_ts": 1567015818,
                    "content_length": 8791271,
                    "student_pid": "102030405060708090000001"
                },
                {
                    "pid": "5d67140850ea66aa1b0ffa76",
                    "media_type": "video",
                    "media_name": "student1_video2.mp4",
                    "media_url": "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1_video2.mp4",
                    "rank_score": 88.33514734096354,
                    "create_ts": 1567015822,
                    "content_length": 5514508,
                    "student_pid": "102030405060708090000001"
                }
            ]
        }
