#!/bin/bash
for i in $(cat cl)
do
	cd $i
	if test $i == "chart"
	then
		go run main.go -chartitle "y=sin(x)" -area -bar -xlabel 5 < sine.d  &
	else 
		go run main.go &
	fi
	cd ..
done
