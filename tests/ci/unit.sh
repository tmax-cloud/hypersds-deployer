#!/bin/bash

go clean -testcache
ginkgo -r -v --skip="\[INCOMPLETE\]" --skipPackage=e2e
