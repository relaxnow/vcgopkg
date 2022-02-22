#!/bin/bash
#Date:22/02/2022
#Author:Boy Baukema
#Purpose:Install Winres and Veracode GO Packager

go install github.com/tc-hib/go-winres@latest
go-winres make
GOOS=windows GOARCH=amd64 go build -o vcgopkg-amd64.exe

#END