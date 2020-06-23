#!/bin/bash
for i in $(cat cl)
do
	cd $i
	if test $i == "chartest"
	then
		./$i -width 1000 -height 1000 -chartitle "y=sin(x)" -area -bar -xlabel 5 < sine.d &
	else 
		./$i &
	fi
	cd ..
done
