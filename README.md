[![Build Status](https://travis-ci.org/RyanJarv/coderun.svg?branch=master)](https://travis-ci.org/RyanJarv/coderun)

This is just an initial concept and will mostly likely not work for you currently.

## Goal
Running code in an isolated environment should be stupid easy (and secure)

## Providers/Resources
### Current
* Docker
  * Python
* Mount
  * AWS Credentials
  
  Mount registers files/directories and optionally allows you to forward the hosts version to the script environment on access


### IceBox
Moved these to ./coderun/icebox until the code is more stable
* Lambda
  * Python
* docker
  * Go
  * NodeJS
  * Ruby
  * Rails

### Docker/Bash/Mount example
```
$ export DOCKER_API_VERSION=1.35 #TODO: figure out why this is needed
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
```

### Lambda/Python example
```
$ export DOCKER_API_VERSION=1.35
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

## Development
The code is meant to be extendable to support easily adding languages/features in the future. Everything is broken down into providers, resoruces, and steps.

A provider introduces a new feature (docker, lambda, file sharing, etc..) which is passed a run environment it then has the option to register any steps and/or attempt to register it's resources. Resources are variations of a feature to support a specific implentation such as a runtime language or a specific host/directory to share. Resources also recieve the run environment and can register additional steps.

Steps are callbacks with some attached info about when they should run which can be one of setup, deploy, run, teardown or a custom stage somewhere inbetween. They can also be registered to run before/after any matching provider/resource/step combination, so you could say run a given callback that has a provider of foo and a resource of bar and matches any step.
