#!/bin/bash
for i in $(cat cl)
do
	cd $i
	go build .
	cd ..
done
