#!/bin/bash

rm -fv ./pkg/flexy/flexy.test
rm -rfv ./bulid/flexy
mkdir -p ./build/flexy
ginkgo build pkg/flexy/
mv -v pkg/flexy/flexy.test ./build/flexy/

#TO run
# ginkgo -v -focus="\[Main\] Flexy" ./build/flexy/flexy.test
## OR,
# ginkgo -v -focus="\[Main\] Flexy" ./pkg/flexy