#!/bin/bash

export SPEECHSDK_ROOT=${PWD}/speechsdk
export CGO_CFLAGS="-I${SPEECHSDK_ROOT}/include/c_api"
export CGO_LDFLAGS="-L${SPEECHSDK_ROOT}/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="${SPEECHSDK_ROOT}/speechsdk/lib/x64:${LD_LIBRARY_PATH}"
export CGO_ENABLED=1

go build
