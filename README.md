[!Build Status](https://travis-ci.org/RyanJarv/coderun.svg?branch=master)](https://travis-ci.org/RyanJarv/coderun)
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
% go run ../../main.go -l info -- ./test.py
INFO2018/02/11 21:27:02 Pulling image: python:3
INFO2018/02/11 21:27:03 Running: [python -m venv .coderun/venv]
WARN2018/02/11 21:27:05 No step deploy registered for provider docker
INFO2018/02/11 21:27:05 Running: [.coderun/venv/bin/python ./test.py]
hello from docker
% go run ../../main.go -l info -p lambda -- ./test.py
INFO2018/02/11 21:28:17 Creating zip file at .coderun/lambda-xvlbzgbaicmrajwwhthc.zip
INFO2018/02/11 21:28:17 Found file: test.py
INFO2018/02/11 21:28:17 Done zipping
START RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e Version: $LATEST
hello from lambda
END RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e
REPORT RequestId: c45aac5a-0fa4-11e8-b6ea-cf2bba77654e  Duration: 0.33 ms   Billed Duration: 100 ms     Memory Size: 128 MB Max Memory Used: 21 MB  
```

