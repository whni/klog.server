conn = Mongo();
db = conn.getDB("klog");
db.dropDatabase();

// institute collection
db.createCollection("institute", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["institute_uid", "institute_name", "address"],
            properties: {
                institute_uid: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z]{1}[a-zA-Z0-9_\-]{5,}+$",
                    minLength: 6,
                    description: "required string (>= 6 length, start with letter)"
                },
                institute_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                address: {
                    bsonType: "object",
                    required: ["city", "state", "country"],
                    description: "required object with country/state/city fields",
                    properties: {
                        street: {
                            bsonType: "string",
                            description: "optional string",
                        },
                        code: {
                            bsonType: "string",
                            minLength: 5,
                            maxLength: 6,
                            description: "optional string",
                        },
                        city: {
                            bsonType: "string",
                            minLength: 2,
                            description: "required string (>= 2 length)"
                        },
                        state: {
                            bsonType: "string",
                            minLength: 2,
                            description: "required string (>= 2 length)"
                        },
                        country: {
                            bsonType: "string",
                            minLength: 2,
                            maxLength: 3,
                            description: "required string (2~3 length)"
                        }
                    }
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.institute.createIndex({"institute_uid": 1}, {unique: true});
db.institute.createIndex({"institute_name": 1}, {unique: true});


// course collection
db.createCollection("course", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["course_uid", "course_name", "course_intro", "teacher_pid", "assistant_pid", "institute_pid"],
            properties: {
                course_uid: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z]{1}[a-zA-Z0-9_\-]{5,}+$",
                    minLength: 6,
                    description: "required string (>= 6 length, start with letter)"
                },
                course_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                course_intro: {
                    bsonType: "string",
                    description: "required string"
                },
                course_targets: {
                    bsonType: ["array"],
                    items: {
                        bsonType: "object",
                        required: ["tag", "desc"],
                        properties: {
                            tag: {
                                bsonType: "string",
                                minLength: 1,
                                description: "required string (>= 1 length)"
                            },
                            desc: {
                                bsonType: "string",
                                minLength: 1,
                                description: "required string (>= 1 length)"
                            }
                        }
                    },
                    description: "optional course target description array"
                },
                teacher_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                assistant_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                institute_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.course.createIndex({"course_uid": 1}, {unique: true});


// teacher collection
db.createCollection("teacher", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["teacher_uid", "teacher_name", "teacher_key", "phone_number", "email", "institute_pid"],
            properties: {
                teacher_uid: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z]{1}[a-zA-Z0-9_\-]{5,}+$",
                    minLength: 6,
                    description: "required string (>= 6 length, start with letter)"
                },
                teacher_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                teacher_key: {
                    bsonType: "string",
                    description: "required string (sha256 hash)"
                },
                phone_number: {
                    bsonType: "string",
                    description: "required string"
                },
                email: {
                    bsonType: "string",
                    description: "required string"
                },
                institute_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.teacher.createIndex({"teacher_uid": 1}, {unique: true});


// student collection
db.createCollection("student", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["student_name", "student_image_name", "student_image_url", "binding_code", "binding_expire"],
            properties: {
                student_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                student_image_name: {
                    bsonType: "string",
                    description: "required string"
                },
                student_image_url: {
                    bsonType: "string",
                    description: "required string"
                },
                binding_code: {
                    bsonType: "string",
                    description: "required string"
                },
                binding_expire: {
                    bsonType: "long",
                    description: "required int64 (unix timestamp)"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});


// relative collection
db.createCollection("relative", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["relative_name", "relative_wxid", "phone_number", "email"],
            properties: {
                relative_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                relative_wxid: {
                    bsonType: "string",
                    minLength: 8,
                    description: "required string (>= 8 length)"
                },
                phone_number: {
                    bsonType: "string",
                    description: "required string"
                },
                email: {
                    bsonType: "string",
                    description: "required string"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.relative.createIndex({"relative_wxid": 1}, {unique: true});


// course_record collection
db.createCollection("course_record", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["student_pid", "course_pid", "target_tag", "record_ts", "is_makeup"],
            properties: {
                student_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                course_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                target_tag: {
                    bsonType: "string",
                    description: "required string - related to course target tag"
                },
                record_ts: {
                    bsonType: "long",
                    description: "required int64 (unix timestamp)"
                },
                is_makeup: {
                    bsonType: "bool",
                    description: "required boolean type to indicate if this is a makeup course record"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});


// course_comment collection
db.createCollection("course_comment", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["course_record_pid", "comment_person_pid", "comment_person_type", "comment_ts", "comment_body"],
            properties: {
                course_record_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                comment_person_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId of the person who gives this comment"
                },
                comment_person_type: {
                    bsonType: "string",
                    enum: ["teacher", "relative", "others"],
                    description: "required person type string: teacher/relative/others"
                },
                comment_ts: {
                    bsonType: "long",
                    description: "required int64 (unix timestamp)"
                },
                comment_body: {
                    bsonType: "string",
                    description: "required comment body string"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});


// cloudmedia collection
db.createCollection("cloudmedia", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["student_pid", "course_record_pid", "media_type", "media_name", "media_url", "rank_score", "media_tags", "create_ts", "content_length"],
            properties: {
                student_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                course_record_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                media_type: {
                    bsonType: "string",
                    enum: ["video", "image", "others"],
                    description: "required string - video/image/others"
                },
                media_name: {
                    bsonType: "string",
                    minLength: 1,
                    description: "required string (media blob name)"
                },
                media_url: {
                    bsonType: "string",
                    minLength: 1,
                    description: "required string (media blob full url)"
                },
                media_tags: {
                    bsonType: ["array"],
                    items: {
                        bsonType: "string",
                        description: "required tag string"
                    },
                    description: "media understanding tags"
                },
                rank_score: {
                    bsonType: "double",
                    description: "required double/float64"
                },
                create_ts: {
                    bsonType: "long",
                    description: "required int64 (unix timestamp)"
                },
                content_length: {
                    bsonType: "long",
                    description: "required int64"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.cloudmedia.createIndex({"media_name": 1}, {unique: true});


// student-relative reference
db.createCollection("student_relative_ref", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["student_pid", "relative_pid", "relationship", "is_main"],
            properties: {
                student_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                relative_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                relationship: {
                    bsonType: "string",
                    description: "required string (e.g., father, mother, uncle, aunt, etc.)"
                },
                is_main: {
                    bsonType: "bool",
                    description: "required boolean type to indicate if this is main relationship"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.student_relative_ref.createIndex( { "student_pid": 1, "relative_pid": 1 }, { unique: true } );

// course-student reference
db.createCollection("student_course_ref", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["course_pid", "student_pid"],
            properties: {
                course_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                },
                student_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.student_course_ref.createIndex( { "course_pid": 1, "student_pid": 1 }, { unique: true } );


// db info
print(`[DB] ${db.getName()} [Collections] ${db.getCollectionNames()}`);

