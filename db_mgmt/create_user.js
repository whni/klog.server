conn = Mongo();
db = conn.getDB("klog");

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
 