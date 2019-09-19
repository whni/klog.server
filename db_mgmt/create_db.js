conn = Mongo();
db = conn.getDB("klog");
db.dropDatabase();

// institute collection
db.createCollection("institutes", {
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
db.institutes.createIndex({"institute_uid": 1}, {unique: true});
db.institutes.createIndex({"institute_name": 1}, {unique: true});


// teacher collection
db.createCollection("teachers", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["teacher_uid", "teacher_name", "teacher_key", "class_name", "phone_number", "email", "institute_pid"],
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
                class_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
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
db.teachers.createIndex({"teacher_uid": 1}, {unique: true});

// student collection
db.createCollection("students", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["student_name", "student_image_name", "student_image_url", "parent_wxid", "parent_name", "phone_number", "email", "binding_code", "binding_expire", "teacher_pid"],
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
                parent_wxid: {
                    bsonType: "string",
                    description: "required string"
                },
                parent_name: {
                    bsonType: "string",
                    description: "required string"
                },
                phone_number: {
                    bsonType: "string",
                    description: "required string"
                },
                email: {
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
                },
                teacher_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
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
            required: ["media_type", "media_name", "media_url", "rank_score", "student_pid", "create_ts", "content_length"],
            properties: {
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
                rank_score: {
                    bsonType: "double",
                    description: "required double/float64"
                },
                student_pid: {
                    bsonType: "objectId",
                    description: "required ObjectId"
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


// db info
print(`[DB] ${db.getName()} [Collections] ${db.getCollectionNames()}`);

