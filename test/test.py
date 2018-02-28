#!/usr/bin/env python
print("hello from python")

import awscli
print(dir(awscli))

def lambda_handler(a, b):
  print("hello from lambda")
  print(dir(awscli))
