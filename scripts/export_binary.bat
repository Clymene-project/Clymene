@echo off
set DOCKERID=build
set RELEASE=1_1_0

mkdir ./out/clymene-%RELEASE%-linux-amd64/
mkdir ./out/clymene-%RELEASE%-windows-amd64/
mkdir ./out/clymene-%RELEASE%-darwin-amd64/
mkdir ./out/clymene-%RELEASE%-linux-s390x/
mkdir ./out/clymene-%RELEASE%-linux-arm64/
mkdir ./out/clymene-%RELEASE%-linux-ppc64le/

docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-amd64 ./out/clymene-%RELEASE%-linux-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-agent-windows-amd64 ./out/clymene-%RELEASE%-windows-amd64/clymene-agent-windows-amd64.exe
docker cp %DOCKERID%:/clymene/out/clymene-agent-darwin-amd64 ./out/clymene-%RELEASE%-darwin-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-s390x ./out/clymene-%RELEASE%-linux-s390x/
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-arm64 ./out/clymene-%RELEASE%-linux-arm64/
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-ppc64le ./out/clymene-%RELEASE%-linux-ppc64le/


docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-amd64 ./out/clymene-%RELEASE%-linux-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-darwin-amd64 ./out/clymene-%RELEASE%-darwin-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-arm64 ./out/clymene-%RELEASE%-linux-arm64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-ppc64le ./out/clymene-%RELEASE%-linux-ppc64le/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-s390x ./out/clymene-%RELEASE%-linux-s390x/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-windows-amd64 ./out/clymene-%RELEASE%-windows-amd64/clymene-ingester-windows-amd64.exe

docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-amd64 ./out/clymene-%RELEASE%-linux-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-darwin-amd64 ./out/clymene-%RELEASE%-darwin-amd64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-arm64 ./out/clymene-%RELEASE%-linux-arm64/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-ppc64le ./out/clymene-%RELEASE%-linux-ppc64le/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-s390x ./out/clymene-%RELEASE%-linux-s390x/
docker cp %DOCKERID%:/clymene/out/clymene-ingester-windows-amd64 ./out/clymene-%RELEASE%-windows-amd64/clymene-ingester-windows-amd64.exe