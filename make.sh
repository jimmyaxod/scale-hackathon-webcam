#!/bin/bash

TAG=$1

# Generate signature
#~/sc/scale signature generate webcam:${TAG} -d signature

# Generate function
~/sc/scale function build -d fn1
