#!/bin/bash
output_dir="$(pwd)/output"
currentdate=$(date +"%Y%m%d")
version="1.0.0-"${currentdate}""
build(){
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "This is a Linux system"
        go build -ldflags "-s -w -X main.version=${version}" -o fileagent main.go
        upx -9 ${output_dir}/fileagent-linux-x86_64
        echo "build success"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "This is a macOS system"
        go build -ldflags "-X main.version=${version}" -o fileagent main.go
        echo "build success"
    elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        echo "This is a Windows system"
        go build -ldflags "-X main.version=${version}" -o fileagent.exe main.go
        echo "build success"
    else
        echo "Unsupported operating system"
    fi
}

buildAll(){
    version=$1
    echo "build version: ${version}"
    mkdir -p ${output_dir}
    echo "build windows"
    go build -ldflags "-s -w -X main.version=${version}" -o ${output_dir}/fileagent-windows-x86_64.exe main.go
    # go build -ldflags "-X main.version=${version}" -o fileagent.exe main.go
    echo "build linux amd64"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=${version}" -o ${output_dir}/fileagent-linux-x86_64 main.go
    upx -9 ${output_dir}/fileagent-linux-x86_64
    # echo "build linux arm64"
    # CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=${version}" -o ${output_dir}/fileagent-linux-aarch64 main.go
    # echo "build macos amd64"
    # CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${version}" -o ${output_dir}/fileagent-darwin-x86_64 main.go
    # echo "build macos arm64"
    # CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${version}"  -o ${output_dir}/fileagent-darwin-aarch64 main.go
    # # 拷贝config.yml
    # cp config-example.yml ${output_dir}/
    echo "build success"
}

buildLinux() {
    version=$1
    echo "build version: ${version}"
    mkdir -p ${output_dir}
    echo "build linux amd64"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags dev -ldflags "-X main.version=${version}" -o ${output_dir}/fileagent-linux-x86_64 main.go
    echo "build success"
}

case "$1" in
"buildAll")
buildAll "$2"
;;
"linux")
    buildLinux "$2"
    ;;
"build")
build
;;
*)
    exit 1
    ;;
esac
