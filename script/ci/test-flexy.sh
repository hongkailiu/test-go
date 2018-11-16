#!/bin/bash

ginkgo -v -skip="\[Main\] Flexy" ./build/flexy/flexy.test
#ginkgo -v -skip="\[Main\] Flexy" ./pkg/flexy/