# go
Golang codebase from Safecility for accessing device streams

## lib
A basic lib for handling universal elements between services and some utility functions

## setup 
Unify how our microservices do their initial setup 
Some utility functions for sql, redis, etc
Secrets some common sense wrappers for gcloud's secretManager

## mqtt
Wrap paho and provide some handling for the things network

## coap

a simple coap broker - mqtt and coap should produce the same stream of SimpleMessages given the same device:
Only transport specific messages should differ and these should be split from other message elements
