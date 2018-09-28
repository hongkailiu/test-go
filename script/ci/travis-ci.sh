#!/bin/bash

rm -rfv build

ginkgo build flexy/
mv flexy/flexy.test build

