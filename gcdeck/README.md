# gcdeck

Interactive deck client.

gcdeck shows slide decks formatted in the [```deck``` markup language](https://github.com/ajstarks/deck/blob/master/README.md)


```gcdeck t.xml``` makes:


![gcdeck](gcdeck.png)


```decksh test.dsh | gcdeck - ```


![gcdeck](gcdeck0.png)

## Keyboard commands

* A, Ctrl-A, ^, 1, Home: first slide
* E, Crtl-E, $, End: last slide
* J, B, Ctrl-B, Ctrl-P, Shift-Space, Shift-Enter: previous slide
* K, F, Ctrl-F, Ctrl-N, Space,       Enter:       previous slide
* G: toggle a grid
* Q, ESC: Quit

## Mouse interactions

* Left Button: next slide
* Right Button: previous slide
* Middle Button: first slide

## Options

```
gcdeck [options] file ("-" for standard input)

Options:

  -page int
    	initial page (default 1)
  -pagesize string
    	pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen (default "Letter")
  -title string
    	slide title
```