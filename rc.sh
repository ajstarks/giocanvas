#!/bin/bash
for i in $(cat cl)
do
	cd $i
	if test $i == "showimage"
	then
		showimage showimage.png &
	else 
		./$i &
	fi
	cd ..
done
