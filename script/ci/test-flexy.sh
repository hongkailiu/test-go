#!/bin/bash

ginkgo -v -skip="\[Main\] Flexy" ./build/flexy.test
#ginkgo -v -skip="\[Main\] Flexy" ./flexy/