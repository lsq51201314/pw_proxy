@echo off
cd %cd%
echo ���ڱ���Windowsƽ̨ ������
go build -ldflags "-s -w" main.go 
echo ������� ������
pause