[![Build Status](https://travis-ci.org/RyanJarv/coderun.svg?branch=master)](https://travis-ci.org/RyanJarv/coderun)

This is just an initial concept and will mostly likely not work for you currently.

## Goal
Running code in an isolated environment should be stupid easy

## Providers
* docker
* lambda

## Languages
* python
* ruby
* nodejs
* go
* bash

## Frameworks
* rails

## Example
```
$ export DOCKER_API_VERSION=1.35 #TODO: figure out why this is needed
$
$ go run ../main.go -- test.sh 
2018/03/03 20:11:53 Settin up MountProvider
{"status":"Pulling from library/ubuntu","id":"latest"}
{"status":"Digest: sha256:e27e9d7f7f28d67aa9e2d7540bdc2b33254b452ee8e60f388875e5b7d9b2b696"}
{"status":"Status: Image is up to date for ubuntu:latest"}
hello from bash
running `cat ~/.aws/credentials` from inside docker

***************************************************
!!! Script is attempting to read ~/.aws/credentials
***************************************************
Is this expected? [yes/no] yes

***super secret keys stored on the host machine***
bye from bash
$ 
$ go run ../../main.go -l info -p lambda -- ./test.py
INFO2018/02/11 21:28:17 Creating zip file at .coderun/lambda-xvlbzgbaicmrajwwhthc.zip
INFO2018/02/11 21:28:17 Found file: test.py
INFO2018/02/11 21:28:17 Done zipping
START RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e Version: $LATEST
hello from lambda
END RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e
REPORT RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e  Duration: 0.33 ms   Billed Duration: 100 ms     Memory Size: 128 MB Max Memory Used: 21 MB  
$
```

