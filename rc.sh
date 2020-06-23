#!/bin/bash
for i in $(cat cl)
do
	cd $i
	./$i &
	cd ..
done
