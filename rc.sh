#!/bin/bash
for i in $(cat cl)
do
	cd $i
	echo -n "$i "
	case $i in
	showfonts)
			./showfonts -text 'Hello Gio' OpenSans*.ttf &
			;;
	showimage)
			./showimage showimage.png &
			;;
	elections)
			./elections nyt-????.d &
			./elections -shape h nyt-19??.d &
			./elections -shape s nyt-19??.d &
			./elections -bgcolor linen -textcolor black -shape l nyt-19??.d &
			./elections -bgcolor linen -textcolor black -shape p nyt-19??.d &
			./elections -bgcolor linen -textcolor black -shape g nyt-19??.d &
			;;
	gcdeck)
			./gcdeck -pagesize 800x500 test.xml &
			;;
	gchart)
			./allcharts 
			;;
	*)
			./$i &
			;;
	esac
	cd ..
done
echo
