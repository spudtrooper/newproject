#!/bin/sh

out=newapp
rm -rf $out
go run main.go -name $out
rm -rf $out