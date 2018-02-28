#! /usr/bin/env node

console.log("Hello from nodejs");

exports.lambda_handler = function(event, context, callback) {
  console.log("Hello from lambda");
}
