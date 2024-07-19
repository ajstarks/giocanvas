#!/bin/bash
for i in $(cat cl)
do
	cd $i
	echo -n "$i "
	case $i in
		showimage)
			./showimage showimage.png & 
			;;
		elections)
			./elections nyt-????.d &
			;;
		gcdeck)
			./gcdeck -pagesize 800x500 .xml &
			;;
		*)
			./$i &
			;;
	esac
	cd ..
done
echo
