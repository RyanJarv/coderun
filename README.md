[![Build Status](https://travis-ci.org/RyanJarv/coderun.svg?branch=master)](https://travis-ci.org/RyanJarv/coderun)
## Update 11/4/2020
Right now thinking this would be interesting to approach from the other direction, going back to a parser in golang, but focusing on integrating it with state machine. Shell's tend to get used for data processing (so I'm told), so having a shell directly on AWS integrating natively with other services seems like it may be useful to someone (maybe even me if I find a job sometime soon). I think I'll start from scratch to get a PoC working here first to further develop the idea, this code however may be useful once I have a better idea where I'm going with this.

I think it makes sense to approach from the parsing side of things due to the need of failing if a statement isn't picked up by our PATH/aliases/overrides/or whatever. Also I don't think you can override things like pipes and redirection operators. Either way though I believe this will be easier for me to understand from the direction of writing a parser than trying to piece together various operations with obscure shell functionality.

[For now I'm using this repo for the PoC](https://github.com/RyanJarv/msh/blob/main/README.md)

## Update 4/8/2020
After taking a break from this, mostly stuck and thinking about where I wanted to take this project next, I believe I need to take a different approach. I still want to support the same high level goals, but I think for now, this should behave more like a package manager focused on running interactive sandboxed containers on Darwin.

To reduce scope from what I already have here I will likely for now be dropping [dockersnitch](https://github.com/RyanJarv/dockersnitch) and lambda, these are both something I would very much like to support (along with many other things) but for now I'm thinking the package manager and config aspects of this is what I need to focus on in the short term. The rest of this section is what I believe this should look like.

At the moment I believe packages will consist of config like files in a directory in the user's path. These config's will behave like executables by referencing coderun on the shebang line, meaning coderun will get executed with the path to the config file as the first argument. The config files may need to contain executable code, but when they do this should be clear and easily auditable.

So from a high level perspective if Dockerfile is for building containers and docker-compose.yml is for dev environment's then coderun.yml (or whatever) would be for end user's.

Since I'm using the word sandbox here I should mention I realize docker isn't often considered a great security boundry. The threat model I'm considering, at least initially, is not unpatched or zero day exploits but instead anything you can always do with normal unix tools as an unprivliged user on a admin's computer. Although kernel exploits and related do seem like less of an issue since Docker on Mac run's in a VM, I am considering that coincidental and don't consider kernel/docker exploits in the scope of this project.

For something similar in behvior on Linux see Jess Fraz's [dockerfiles](https://github.com/jessfraz/dockerfiles) repo. Her project is also in part the inspiration for the next version of coderun.

Everything below here is about the original coderun project (will keep this on the [coderun-v1](https://github.com/RyanJarv/coderun/tree/coderun-v1) branch).

## Goal
Running code in an isolated environment should be stupid easy (and secure)

NOTE: This is just an initial concept and will mostly likely not work for you currently.

## Providers/Resources
### Current
* Docker
  * Bash
  
### IceBox
Moved these to ./coderun/icebox until the code is more stable
* Lambda
  * Python
* docker
  * Go
  * NodeJS
  * Ruby
  * Rails
* Mount (Prompts on shared file access)
  * AWS Credentials
* Snitch (Prompts on new connection attempt)
  * Docker

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

## Ideas
* Watch for new listening ports and ask to forward them to the host
  * Can docker do this after a container is started?
* Make a few repos as examples of well hidden backdoors this tool can protect against
* Some of this might make sense as integrations to other tools like bundler, virtualenv, npm, kitchen, terraform etc..
* Automatic and opinionated dependency management

## Development
The code is meant to be extendable to support easily adding languages/features in the future. Everything is broken down into providers, resoruces, and steps.

A provider introduces a new feature (docker, lambda, file sharing, etc..) which is passed a run environment it then has the option to register any steps and/or attempt to register it's resources. Resources are variations of a feature to support a specific implentation such as a runtime language or a specific host/directory to share. Resources also recieve the run environment and can register additional steps.

Steps are callbacks with some attached info about when they should run which can be one of setup, deploy, run, teardown or a custom stage somewhere inbetween. They can also be registered to run before/after any matching provider/resource/step combination, so you could say run a given callback that has a provider of foo and a resource of bar and matches any step.

