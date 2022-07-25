#!/bin/bash

make build-all-platforms

mkdir ../out
mkdir ../out/clymene-v2_1_1-linux-amd64/
mkdir ../out/clymene-v2_1_1-windows-amd64/
mkdir ../out/clymene-v2_1_1-darwin-amd64/
mkdir ../out/clymene-v2_1_1-linux-s390x/
mkdir ../out/clymene-v2_1_1-linux-arm64/
mkdir ../out/clymene-v2_1_1-linux-ppc64le/

mv ../out/clymene-agent-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/clymene-agent-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/clymene-agent-windows-amd64.exe
mv ../out/clymene-agent-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/clymene-agent-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/clymene-agent-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/clymene-agent-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/


mv ../out/clymene-ingester-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/clymene-ingester-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/clymene-ingester-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/clymene-ingester-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/
mv ../out/clymene-ingester-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/clymene-ingester-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/clymene-ingester-windows-amd64.exe

mv ../out/clymene-gateway-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/clymene-gateway-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/clymene-gateway-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/clymene-gateway-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/
mv ../out/clymene-gateway-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/clymene-gateway-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/clymene-gateway-windows-amd64.exe

mv ../out/clymene-promtail-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/clymene-promtail-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/clymene-promtail-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/clymene-promtail-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/
mv ../out/clymene-promtail-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/clymene-promtail-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/clymene-promtail-windows-amd64.exe

mv ../out/promtail-ingester-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/promtail-ingester-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/promtail-ingester-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/promtail-ingester-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/
mv ../out/promtail-ingester-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/promtail-ingester-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/promtail-ingester-windows-amd64.exe

mv ../out/promtail-gateway-linux-amd64 ../out/clymene-v2_1_1-linux-amd64/
mv ../out/promtail-gateway-darwin-amd64 ../out/clymene-v2_1_1-darwin-amd64/
mv ../out/promtail-gateway-linux-arm64 ../out/clymene-v2_1_1-linux-arm64/
mv ../out/promtail-gateway-linux-ppc64le ../out/clymene-v2_1_1-linux-ppc64le/
mv ../out/promtail-gateway-linux-s390x ../out/clymene-v2_1_1-linux-s390x/
mv ../out/promtail-gateway-windows-amd64 ../out/clymene-v2_1_1-windows-amd64/promtail-gateway-windows-amd64.exe