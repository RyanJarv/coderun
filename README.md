This is just an initial concept and will mostly likely not work for you currently.

## Goal
Running code in an isolated environment should be stupid easy

## Languages
* python
* ruby
* nodejs
* go
* bash

## Example
```
% go run main.go -- test.rb 
2018/02/03 04:14:13 [/usr/local/bin/docker pull ruby:2.1]
2018/02/03 04:14:13 Running command and waiting for it to finish...
2018/02/03 04:14:13 output: 2.1: Pulling from library/ruby
Digest: sha256:568664cfb53cd74147590cc86e356c4e352e00e28a0011ea1443e8664ca5bad5
Status: Image is up to date for ruby:2.1
2018/02/03 04:14:13 [/usr/local/bin/docker run -t --rm --name my-running-script1 -v /Users/rgerstenkorn/Code/coderun:/usr/src/myapp -w /usr/src/myapp ruby:2.1 bundler install --path .coderun/vendor/bundle]
2018/02/03 04:14:13 Running command and waiting for it to finish...
2018/02/03 04:14:15 output: Warning: the running version of Bundler (1.15.1) is older than the version that created the lockfile (1.16.1). We suggest you upgrade to the latest version of Bundler by running `gem install bundler`.
Using rspec-support 3.7.1
Using diff-lcs 1.3
Using bundler 1.15.1
Using rspec-core 3.7.1
Using rspec-expectations 3.7.0
Using rspec-mocks 3.7.0
Using rspec 3.7.0
Bundle complete! 1 Gemfile dependency, 7 gems now installed.
Bundled gems are installed into ./.coderun/vendor/bundle.
2018/02/03 04:14:15 randString: xvlbzgbaicmrajw
2018/02/03 04:14:15 [/usr/local/bin/docker run -t --rm --name coderun-xvlbzgbaicmrajw -v /Users/rgerstenkorn/Code/coderun:/usr/src/myapp -w /usr/src/myapp ruby:2.1 ruby test.rb]
2018/02/03 04:14:15 Running command and waiting for it to finish...
2018/02/03 04:14:16 output: hello world from ruby
```

