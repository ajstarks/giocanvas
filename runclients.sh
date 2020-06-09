#!/bin/bash
for i in chart concentric confetti eclipse hello lines mondrian play rl sunearth
do
	cd $i
	if test $i == "chart"
	then
		go run main.go -chartitle "y=sin(x)" -area -bar -xlabel 5 < sine.d&
	else 
		go run main.go&
	fi
	cd ..
done
