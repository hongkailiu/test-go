#!/bin/bash

rm -fv ./pkg/flexy/flexy.test
rm -rfv ./bulid/
mkdir -p ./build/
ginkgo build pkg/flexy/
mv -v pkg/flexy/flexy.test ./build/

#TO run
# ginkgo -v -focus="\[Main\] Flexy" ./build/flexy.test
## OR,
# ginkgo -v -focus="\[Main\] Flexy" ./pkg/flexy