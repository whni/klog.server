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
                    maxLength: 32,
                    description: "required string (md5 hash)"
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
            required: ["student_uid", "student_name", "parent_wxid", "parent_name", "phone_number", "email", "teacher_pid"],
            properties: {
                student_uid: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z]{1}[a-zA-Z0-9_\-]{5,}+$",
                    minLength: 6,
                    description: "required string (>= 6 length, start with letter)"
                },
                student_name: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                parent_wxid: {
                    bsonType: "string",
                    minLength: 2,
                    description: "required string (>= 2 length)"
                },
                parent_name: {
                    bsonType: "string",
                    minLength: 2,
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
db.students.createIndex({"student_uid": 1}, {unique: true});

// db info
print(`[DB] ${db.getName()} [Collections] ${db.getCollectionNames()}`);


// init data
db.institutes.insertMany(
    [
        {
            _id: ObjectId("102030405060708090000001"),
            institute_uid: "uid-usa-0001",
            institute_name: "Institute 1",
            address: {
                street: "180 Elm Ct",
                code: "94086",
                city: "Sunnyvale",
                state: "CA",
                country: "USA"
            }
        },
        {
            _id: ObjectId("102030405060708090000002"),
            institute_uid: "uid-usa-0002",
            institute_name: "Institute 2",
            address: {
                street: "Valley Green 6",
                code: "95014",
                city: "Cupertino",
                state: "CA",
                country: "USA"
            }
        }
    ]
);

db.teachers.insertMany(
    [
        {
            _id: ObjectId("102030405060708090000001"),
            teacher_uid: "uid-usa-1001",
            teacher_name: "Nicole Taylor",
            teacher_key: "no_key",
            class_name: "GoldenEye",
            phone_number: "123-456-9876",
            email: "nigoo@klog.com",
            institute_pid: ObjectId("102030405060708090000001")
        },
        {
            _id: ObjectId("102030405060708090000002"),
            teacher_uid: "uid-usa-1002",
            teacher_name: "Wayne Grace",
            teacher_key: "no_key",
            class_name: "FastWind",
            phone_number: "123-456-9876",
            email: "wayne@klog.com",
            institute_pid: ObjectId("102030405060708090000001")
        },
        {
            _id: ObjectId("102030405060708090000003"),
            teacher_uid: "uid-usa-1003",
            teacher_name: "Fantasy God",
            teacher_key: "no_key",
            class_name: "CloudTop",
            phone_number: "000-111-2222",
            email: "fanfan@klog.com",
            institute_pid: ObjectId("102030405060708090000002")
        },
        {
            _id: ObjectId("102030405060708090000004"),
            teacher_uid: "uid-usa-1004",
            teacher_name: "Summer Season",
            teacher_key: "no_key",
            class_name: "UnderWorld",
            phone_number: "619-763-1020",
            email: "summer@klog.com",
            institute_pid: ObjectId("102030405060708090000002")
        }
    ]
);




db.students.insertMany(
    [
        {
            _id: ObjectId("102030405060708090000001"),
            student_uid: "uid-usa-1001",
            student_name: "Thomas Hu",
            parent_wxid: "wxid-0123456789",
            parent_name: "Ed Sheeran",
            phone_number: "777-888-9999",
            email: "ed_sh@apple.com",
            teacher_pid: ObjectId("102030405060708090000001")
        },
        {
            _id: ObjectId("102030405060708090000002"),
            student_uid: "uid-usa-1002",
            student_name: "Bruce Wang",
            parent_wxid: "wxid-0123456789",
            parent_name: "Madison Beer",
            phone_number: "777-888-9999",
            email: "beer@google.com",
            teacher_pid: ObjectId("102030405060708090000002")
        },
        {
            _id: ObjectId("102030405060708090000003"),
            student_uid: "uid-usa-1003",
            student_name: "Tiffiny Shawn",
            parent_wxid: "wxid-0123456789",
            parent_name: "Skylar Grey",
            phone_number: "777-888-9999",
            email: "skylar@facebook.com",
            teacher_pid: ObjectId("102030405060708090000003")
        },
        {
            _id: ObjectId("102030405060708090000004"),
            student_uid: "uid-usa-1004",
            student_name: "Gintama Y.",
            parent_wxid: "wxid-0123456789",
            parent_name: "Autumn Mendes",
            phone_number: "777-888-9999",
            email: "autumn@xxx.com",
            teacher_pid: ObjectId("102030405060708090000004")
        }
    ]
);
