# shapesketch

![shapesketch](default.png)
![shapesketch-dark](dark-curve.png)

![line](line.png)
![circle](circle.png)
![square](square.png)
![rect](rect.png)
![curve](curve.png)
![ellipse](ellipse.png)
![arc](arc1.png)

## Pointer controls

* Primary pointer press: define begin point
* Secondary pointer press: define ending point
* Tertiary pointer press: show the decksh spec
* Move: define current point

## Keyboard controls

### Adjustments

* left, right, up, down arrow keys: adjust begin point 
* Ctrl + left, up, down, arrow keys: adjust end point

### Shapes

* A: arc
* B: quadradic Bezier
* C: circle
* E: ellipse
* L: line
* R: rectangle
* S: square

### Info

* D: show the decksh spec
* G: toggle coordinate grid
* Q, ESC: quit



## Command-line options
```
  -begincolor string
    	begin coordinate color (default "green")
  -bgcolor string
    	background color (default "white")
  -csize float
    	coordinate size (default 1.25)
  -currentcolor string
    	current coordinate color (default "gray")
  -endcolor string
    	end coordinate color (default "red")
  -height int
    	canvas height (default 1000)
  -lsize float
    	line size (default 1)
  -precision int
    	coordinate precision
  -shapecolor string
    	curve color (default "#22222255")
  -ssize float
      step size for coordinate nudging (default 0.5)
  -textcolor string
    	text color (default "black")
  -tsize float
    	text size (default 2.5)
  -width int
    	canvas width (default 1000)
```
