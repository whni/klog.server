CREATE USER IF NOT EXISTS 'vloguser'@'%' IDENTIFIED BY 'vloguser';
GRANT ALL PRIVILEGES ON vlog.* TO 'vloguser'@'%';