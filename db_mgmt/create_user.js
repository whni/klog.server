conn = Mongo();
db = conn.getDB("klog");

if (db.getUser("klog_user") == null) {
    db.createUser(
        {
        user: "klog_user",
            pwd: "klog_pwd",
            roles: [
                {
                    role: "readWrite",
                    db: "klog"
                }
            ]
        }
    )
}