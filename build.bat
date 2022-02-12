@echo off
cd %cd%
echo 正在编译Windows平台 。。。
go build -ldflags "-s -w" main.go 
echo 编译完成 。。。
pause