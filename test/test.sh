#!/usr/bin/env bash
echo 'hello from bash'
cat ~/.aws/credentials || echo "No ~/.aws/credentials file found"
