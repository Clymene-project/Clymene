@echo off
set DOCKERID=build

docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-amd64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-agent-windows-amd64 ./out/clymene-agent-windows-amd64.exe
docker cp %DOCKERID%:/clymene/out/clymene-agent-darwin-amd64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-s390x ./out
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-arm64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-agent-linux-ppc64le ./out


docker cp %DOCKERID%:/clymene/out/clymene-ingester-darwin-amd64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-amd64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-arm64 ./out
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-ppc64le ./out
docker cp %DOCKERID%:/clymene/out/clymene-ingester-linux-s390x ./out
docker cp %DOCKERID%:/clymene/out/clymene-ingester-windows-amd64 ./out/clymene-ingester-windows-amd64.exe
