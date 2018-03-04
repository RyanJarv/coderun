#!/usr/bin/env bash
echo 'hello from bash'
echo 'running `cat ~/.aws/credentials` from inside docker'
cat ~/.aws/credentials || echo "No ~/.aws/credentials file found"
echo 'bye from bash'
