[![Build Status](https://travis-ci.org/RyanJarv/coderun.svg?branch=master)](https://travis-ci.org/RyanJarv/coderun)

This is just an initial concept and will mostly likely not work for you currently.

## Goal
Running code in an isolated environment should be stupid easy (and secure)

## Providers/Resources
### Current
* Docker
  * Bash
* Mount (Prompts on shared file access)
  * AWS Credentials
* Snitch (Prompts on new connection attempt)
  * Docker
  


### IceBox
Moved these to ./coderun/icebox until the code is more stable
* Lambda
  * Python
* docker
  * Go
  * NodeJS
  * Ruby
  * Rails

### Bash/Mount/Dockersnitch example
```
$ go run ../main.go 
» ./test.sh 
{"status":"Pulling from library/bash","id":"latest"}
{"status":"Digest: sha256:717f5f1e5f15624166a6abaa8f5c99a5f812c379f3a5f1a31db1dd7206ef9107"}
{"status":"Status: Image is up to date for bash:latest"}
hello from bash
running wget -O /tmp/wgetout http://google.com
Connecting to google.com (172.217.0.142:80)
Allow connection from 172.217.0.142? [w/b] w
Connecting to www.google.com (74.125.136.105:80)
Allow connection from 74.125.136.105? [w/b] w
wgetout              100% |*******************************| 12499   0:00:00 ETA
running `cat ~/.aws/credentials` from inside docker
!!! Script is attempting to read ~/.aws/credentials, is this expected? [yes/no] yes
bye from bash
» ^C
$ 
```

### Lambda/Python example
```
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
