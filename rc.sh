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
			;;
	gcdeck)
			./gcdeck -pagesize 800x500 test.xml &
			;;
	gchart)
			./gchart -bar -top 80 -bottom 60 -xlabel 10 cos.d &
			;;
	*)
			./$i &
			;;
	esac
	cd ..
done
echo
