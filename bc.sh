#!/bin/bash
for i in $(cat cl)
do
	cd $i
	echo -n "$i "
	go build $* -ldflags="-s -w" .
	cd ..
done
echo
