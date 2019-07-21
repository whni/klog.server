CREATE USER IF NOT EXISTS 'kloguser'@'%' IDENTIFIED BY 'kloguser';
GRANT ALL PRIVILEGES ON klog_business.* TO 'kloguser'@'%';