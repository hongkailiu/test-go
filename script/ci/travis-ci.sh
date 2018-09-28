#!/bin/bash

rm -rfv build
mkdir build
ginkgo build flexy/
mv flexy/flexy.test build/

