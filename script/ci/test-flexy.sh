#!/bin/bash

#ginkgo -v -skip="\[Main\] Flexy" ./build/flexy.test
pwd
ls -al ./
ls -al build
ginkgo -v -skip="\[Main\] Flexy" ./flexy/