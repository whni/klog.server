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
                    minLength: 5,
                    description: "required string (>= 5 length)"
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
            required: ["teacher_uid", "teacher_name", "class_name", "phone_number", "email", "institute_uid"],
            properties: {
                teacher_uid: {
                    bsonType: "string",
                    minLength: 5,
                    description: "required string (>= 5 length)"
                },
                teacher_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                class_name: {
                    bsontype: "string",
                    minlength: 2,
                    description: "required string (>= 2 length)"
                },
                phone_number: {
                    bsontype: "string",
                    description: "required string"
                },
                email: {
                    bsontype: "string",
                    description: "required string"
                },
                institute_uid: {
                    bsonType: "string",
                    minLength: 5,
                    description: "required string (>= 5 length)"
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
            required: ["student_uid", "student_name", "parent_wxid", "parent_name", "phone_number", "email", "teacher_uid", "institute_uid"],
            properties: {
                student_uid: {
                    bsonType: "string",
                    minLength: 5,
                    description: "required string (>= 5 length)"
                },
                student_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                parent_wxid: {
                    bsontype: "string",
                    minlength: 2,
                    description: "required string (>= 2 length)"
                },
                parent_name: {
                    bsontype: "string",
                    minlength: 2,
                    description: "required string"
                },
                phone_number: {
                    bsontype: "string",
                    description: "required string"
                },
                email: {
                    bsontype: "string",
                    description: "required string"
                },
                teacher_uid: {
                    bsonType: "string",
                    minLength: 5,
                    description: "required string (>= 5 length)"
                },
                institute_uid: {
                    bsonType: "string",
                    minLength: 5,
                    description: "required string (>= 5 length)"
                }
            }
        }
    },
    validationLevel: "strict",
    validationAction: "error"
});
db.students.createIndex({"student_uid": 1}, {unique: true});

// db info
print(`[DB] ${db.getName()} [Collections] ${db.getCollectionNames()}`);


// init data
db.institutes.insertMany(
    [
        {
            institute_uid: "ins-10001",
            institute_name: "test ins 1",
            address: {
                country: "USA",
                state: "CA",
                city: "Sunnyvale",
                street: "Valley Green 6"
            }
        },
        {
            institute_uid: "ins-10002",
            institute_name: "test ins 2",
            address: {
                country: "USA",
                state: "CA",
                city: "Sunnyvale",
                street: "Valley Green 6"
            }
        }
    ]
);
