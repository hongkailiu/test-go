#!/bin/bash

rm -fv flexy/flexy.test
ginkgo build flexy/
mv flexy/flexy.test build

#TO run
# ginkgo -v -focus="\[Main\] Flexy" ./build/flexy.test
## OR,
# ginkgo -v -focus="\[Main\] Flexy" ./flexy