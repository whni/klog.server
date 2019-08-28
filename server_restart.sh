#!/bin/sh

sudo killall klog.server

go build .

sleep 1

sudo ./klog.server &