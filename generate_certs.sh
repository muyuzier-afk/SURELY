#!/bin/bash

# 生成 TLS 证书
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"

echo "TLS certificates generated successfully!"
echo "Files created:"
echo "- server.key (private key)"
echo "- server.crt (certificate)"
