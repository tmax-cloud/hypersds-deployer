#!/bin/bash

go clean -testcache
ginkgo -r -v --skip="\[E2e\]|\[INCOMPLETE\]"
