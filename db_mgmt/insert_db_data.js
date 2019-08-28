conn = Mongo();
db = conn.getDB("klog");

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
            teacher_name: "倪炜恒",
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
            student_name: "Thomas Hu",
            student_image_url: "",
            parent_wxid: "wxid-0123456789",
            parent_name: "Ed Sheeran",
            phone_number: "777-888-9999",
            email: "ed_sh@apple.com",
            binding_code: "",
            binding_expire: NumberLong(0),
            teacher_pid: ObjectId("102030405060708090000001")
        },
        {
            _id: ObjectId("102030405060708090000002"),
            student_name: "Bruce Wang",
            student_image_url: "",
            parent_wxid: "wxid-0123456789",
            parent_name: "Madison Beer",
            phone_number: "777-888-9999",
            email: "beer@google.com",
            binding_code: "",
            binding_expire: NumberLong(0),
            teacher_pid: ObjectId("102030405060708090000002")
        },
        {
            _id: ObjectId("102030405060708090000003"),
            student_name: "Tiffiny Shawn",
            student_image_url: "",
            parent_wxid: "wxid-0123456789",
            parent_name: "Skylar Grey",
            phone_number: "777-888-9999",
            email: "skylar@facebook.com",
            binding_code: "",
            binding_expire: NumberLong(0),
            teacher_pid: ObjectId("102030405060708090000003")
        },
        {
            _id: ObjectId("102030405060708090000004"),
            student_name: "Gintama Y.",
            student_image_url: "",
            parent_wxid: "wxid-0123456789",
            parent_name: "Autumn Mendes",
            phone_number: "777-888-9999",
            email: "autumn@xxx.com",
            binding_code: "",
            binding_expire: NumberLong(0),
            teacher_pid: ObjectId("102030405060708090000004")
        },
        {
            _id: ObjectId("1020304050607080900000FF"),
            student_name: "Baby Cute",
            student_image_url: "https://klogresourcediag.blob.core.windows.net/klog-cloud-media/student1.jpg",
            parent_wxid: "orgQa44wYyOpdShmXAsHtSfjMjeQ",
            parent_name: "Bruce Wayne",
            phone_number: "619-763-4183",
            email: "brucexxx@klog.com",
            binding_code: "",
            binding_expire: NumberLong(0),
            teacher_pid: ObjectId("102030405060708090000004")
        }
    ]
);