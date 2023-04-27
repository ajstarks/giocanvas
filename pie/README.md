# pie chart

Makes pie charts reading from files in this format:
one line per item, fields (name, value, color; tab-separated)
lines beginning with '#' are the title.

For example:

```
# Desktop Browser Market Share 2021-09
Chrome	67.17	red
Edge	9.33	green
Firefox	7.87	orange
Safari	9.63	blue
Other	5.99	gray
```

If no files a specified, embedded data is shown.  The command line options:

```
  -duration duration
    	animation interval (default 1s)
  -height int
    	canvas height (default 1000)
  -width int
    	canvas width (default 1000)
```

![pie](pie.png)
