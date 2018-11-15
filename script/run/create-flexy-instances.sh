#!/bin/bash


ginkgo -v -focus="\[Main\] Flexy" ./build/flexy.test
## OR,
# ginkgo -v -focus="\[Main\] Flexy" ./pkg/flexy