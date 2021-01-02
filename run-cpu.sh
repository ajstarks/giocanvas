#!/bin/bash
for i in $(cat cl)
do
	echo $i
	cd $i
	GIORENDERER=forcecompute go run .
	cd ..
done
