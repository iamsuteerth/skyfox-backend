#!/bin/bash
set -e

# Debug: List directories to verify they exist
echo "Checking directories:"
ls -la /var
ls -la /var/log
ls -la /var/run

# Create directories if they don't exist (as a fallback)
mkdir -p /var/log/supervisor
mkdir -p /var/run/supervisor
mkdir -p /var/log/nginx
mkdir -p /var/run/nginx
mkdir -p /tmp/nginx

# Replace environment variables in nginx config
envsubst '${API_GATEWAY_KEY}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Start all services via supervisord
exec supervisord -c /etc/supervisord.conf