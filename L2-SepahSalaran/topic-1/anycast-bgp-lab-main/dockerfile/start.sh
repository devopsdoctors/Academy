#!/bin/bash
ip addr add 192.168.10.100/32 dev eth0
nginx -g 'daemon off;
