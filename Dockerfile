from golang:alpine

expose 8080
run apk add bash ca-certificates git gcc g++ libc-dev

run mkdir hackatonEcooltra
workdir /hackatonEcooltra

copy . /hackatonEcooltra
run go build .

cmd ["./hackatonEcooltra"] 
