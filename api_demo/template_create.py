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
api_url = "{}/api/0/config/template".format(host)
method = HTTPMethod.POST
template_params = [
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 5,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": [
            {
                "clip_duration": 5,
                "clip_sequence": 1
            }
        ],
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food_1.mp3",
        "template_name": "001_Food_Delightful_short",
        "template_tags": [
            "happy",
            "kids",
            "family",
            "food"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 6,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 5,
                "clip_sequence": 4
            },
            {
                "clip_duration": 5,
                "clip_sequence": 5
            },
            {
                "clip_duration": 4,
                "clip_sequence": 6
            }
        ],
        "template_clip_time_json": [
            {
                "clip_duration": 5,
                "clip_sequence": 1
            }
        ],
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food_3.mp3",
        "template_name": "001_Food_Delightful_long",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 8,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 5,
                "clip_sequence": 4
            },
            {
                "clip_duration": 5,
                "clip_sequence": 5
            },
            {
                "clip_duration": 4,
                "clip_sequence": 6
            },
            {
                "clip_duration": 4,
                "clip_sequence": 7
            },
            {
                "clip_duration": 5,
                "clip_sequence": 8
            }
        ],
        "template_clip_time_json": [
            {
                "clip_duration": 5,
                "clip_sequence": 1
            }
        ],
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_food_2.mp3",
        "template_name": "001_Food_Delightful_long",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            },
            {
                "clip_duration": 6,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 6,
                "clip_sequence": 4
            },
            {
                "clip_duration": 4,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "Vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_StreamLive.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Selfie_1.mp3",
        "template_name": "002_Selfie_LiveStreaming",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 2,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            },
            {
                "clip_duration": 2,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": [
            {
                "1": "girl_selfie_1.mp4"
            }
        ],
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_StreamLive.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Selfie_1.mp3",
        "template_name": "002_Selfie_GirlSelfiea",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 6,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 2
            },
            {
                "clip_duration": 1,
                "clip_sequence": 4
            },
            {
                "clip_duration": 1,
                "clip_sequence": 6
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": [
            {
                "1": "boy_selfie_1.mp4"
            },
            {
                "3": "boy_selfie_2.mp4"
            },
            {
                "5": "boy_selfie_3.mp4"
            }
        ],
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_selfie_3_boy.mp3",
        "template_name": "002_Selfie_BoySelfiea",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 3,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 3,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_pet_1.mp3",
        "template_name": "003_Pet_HappyMusic",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 2,
                "clip_sequence": 3
            },
            {
                "clip_duration": 3,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_pet_2.mp3",
        "template_name": "003_Pet_CuteMusic",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 1
            },
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 2,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            },
            {
                "clip_duration": 1,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_pet_3_dog.mp3",
        "template_name": "003_Pet_Dog",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 1
            },
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 2,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            },
            {
                "clip_duration": 1,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_pet_4_cat.mp3",
        "template_name": "003_Pet_CatMusic",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 2,
        "template_clip_number_total": 3,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 1,
                "clip_sequence": 3
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": [
            {
                "1": "pet_cat_opening.mp4"
            }
        ],
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_pet_5_cat.mp3",
        "template_name": "003_Pet_CatMovie",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 3,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            },
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 6,
                "clip_sequence": 3
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_decoration.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_1.mp3",
        "template_name": "004_Kid_Music_cute",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_decoration.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_2.mp3",
        "template_name": "004_Kid_Music_disco",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 1
            },
            {
                "clip_duration": 3,
                "clip_sequence": 2
            },
            {
                "clip_duration": 6,
                "clip_sequence": 3
            },
            {
                "clip_duration": 3,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Sun_Moon.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_3.mp3",
        "template_name": "004_Kid_Music_quiet",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 3,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            },
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 5,
                "clip_sequence": 3
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_decoration.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_5.mp3",
        "template_name": "004_Kid_Music_birthday_short",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 6,
        "template_clip_number_total": 6,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 2,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            },
            {
                "clip_duration": 2,
                "clip_sequence": 5
            },
            {
                "clip_duration": 3,
                "clip_sequence": 6
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_decoration.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_5.mp3",
        "template_name": "004_Kid_Music_birthday_long",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 1
            },
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 6,
                "clip_sequence": 3
            },
            {
                "clip_duration": 6,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_TextUnderline_kid.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_kid_zimu.mp3",
        "template_name": "004_Kid_Music_text_underline",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 5,
                "clip_sequence": 1
            },
            {
                "clip_duration": 5,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 6,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_TextUnderline_Travel.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Travel_1.mp3",
        "template_name": "005_Vacation_Travel_textUnderline_1",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            },
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 3,
                "clip_sequence": 3
            },
            {
                "clip_duration": 4,
                "clip_sequence": 4
            },
            {
                "clip_duration": 5,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_TextUnderline_Travel2.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Travel_2.mp3",
        "template_name": "005_Vacation_Travel_textUnderline_2",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 3,
                "clip_sequence": 1
            },
            {
                "clip_duration": 3,
                "clip_sequence": 2
            },
            {
                "clip_duration": 3,
                "clip_sequence": 3
            },
            {
                "clip_duration": 3,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Travel_3.mp3",
        "template_name": "005_Vacation_Travel_Music",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 6,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 6,
                "clip_sequence": 4
            },
            {
                "clip_duration": 4,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            }
        ],
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_travel_world_long.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Travel_4.mp3",
        "template_name": "005_Vacation_Travel_Music_2",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 6,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 6,
                "clip_sequence": 4
            },
            {
                "clip_duration": 4,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            }
        ],
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_travel_world_short.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Travel_4.mp3",
        "template_name": "005_Vacation_Travel_Happy",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 5,
        "template_clip_number_total": 5,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 2,
                "clip_sequence": 2
            },
            {
                "clip_duration": 3,
                "clip_sequence": 3
            },
            {
                "clip_duration": 3,
                "clip_sequence": 4
            },
            {
                "clip_duration": 5,
                "clip_sequence": 5
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Natural_View.mp3",
        "template_name": "006_Natural_View_music",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 6,
                "clip_sequence": 1
            },
            {
                "clip_duration": 5,
                "clip_sequence": 2
            },
            {
                "clip_duration": 4,
                "clip_sequence": 3
            },
            {
                "clip_duration": 5,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_city_zimu_1.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_city_zimu_1.mp3",
        "template_name": "007_City_textUnderline",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 6,
                "clip_sequence": 1
            },
            {
                "clip_duration": 6,
                "clip_sequence": 2
            },
            {
                "clip_duration": 6,
                "clip_sequence": 3
            },
            {
                "clip_duration": 2,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_city_zimu_2.json",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_city_zimu_2.mp3",
        "template_name": "007_City_textUnderline_2",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 4,
        "template_clip_number_total": 4,
        "template_clip_time_content": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            },
            {
                "clip_duration": 3,
                "clip_sequence": 2
            },
            {
                "clip_duration": 3,
                "clip_sequence": 3
            },
            {
                "clip_duration": 4,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": "null",
        "template_filter": "vanilla",
        "template_json": "https://klogresourcediag159.blob.core.windows.net/story-template/template_decoration.zip",
        "template_mp4_movie": "null",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_Christmas.mp3",
        "template_name": "008_Festival_Christmas",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 3,
        "template_clip_time_content": [
            {
                "clip_duration": 4,
                "clip_sequence": 2
            },
            {
                "clip_duration": 3,
                "clip_sequence": 3
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": [
            {
                "clip_duration": 4,
                "clip_sequence": 1
            }
        ],
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "https://klogresourcediag159.blob.core.windows.net/story-template/template_thanksgiving_opening.mp4",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_thanksgiving.mp3",
        "template_name": "008_Festival_Thanksgiving",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    },
    {
        "template_clip_number_needed": 3,
        "template_clip_number_total": 3,
        "template_clip_time_content": [
            {
                "clip_duration": 6,
                "clip_sequence": 2
            },
            {
                "clip_duration": 6,
                "clip_sequence": 3
            },
            {
                "clip_duration": 4,
                "clip_sequence": 4
            }
        ],
        "template_clip_time_json": "null",
        "template_clip_time_mp4": [
            {
                "clip_duration": 2,
                "clip_sequence": 1
            }
        ],
        "template_filter": "vanilla",
        "template_json": "null",
        "template_mp4_movie": "https://klogresourcediag159.blob.core.windows.net/story-template/template_general_festival_opening.mp4",
        "template_music": "https://klogresourcediag159.blob.core.windows.net/story-template/template_festival_wedding.mp3",
        "template_name": "008_Festival_general",
        "template_tags": [
            "happy",
            "kids",
            "family"
        ]
    }
]

for params in template_params:
    print("[create template {}]".format(params["template_name"]))
    http_req = HTTPRequest(api_url, method, params)
    http_req.send()
    http_req.print_resp()

# Query template
# url + method
api_url = "{}/api/0/config/template?pid=all&fkey=template_name&fid=008_Festival_general".format(host)
method = HTTPMethod.GET
params = {}
print("[Query template]")
http_req = HTTPRequest(api_url, method, params)
http_req.send()
http_req.print_resp()