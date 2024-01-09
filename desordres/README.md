# decordres -- concentric squares in the manner of Des Ordres by Vera Molnár

![default](default.png)

```
desordres -width=500 -height=500 
```

![hot](hot.png)

```
desordres -width=500 -height=500 -maxlw=0.5 -bgcolor=black -color='20:60'
```



## interaction

* "+", Primary Mouse, Right or Up Arrow: increase tiles/row
* "-", Secondary Mouse, Left or Down Arrow: decrease tiles/row
* "P": randonly select a palette
* "Q", Esc: Quit
* Home: minimum tile/row: 1
* End: maximum tile/row: 20

Pressing any other key will create a new set of tiles

## options

```
Option      Default    Description
.....................................................
-help       false      show usage
-width      1000       canvas width
-height     1000       canvas height
-tiles      10         number of tiles/row
-maxlw      1          maximim line thickness
-bgcolor    white      background color
-color      gray       color name, h1:h2, or palette:
                       2-bit-demichrome               [#211e20 #555568 #a0a08b #e9efec]
                       hollow                         [#0f0f1b #565a75 #c6b7be #fafbf6]
                       ice-cream-gb                   [#7c3f58 #eb6b6f #f9a875 #fff6d3]
                       ayy4                           [#00303b #ff7777 #ffce96 #f1f2da]
                       arq4                           [#ffffff #6772a9 #3a3277 #000000]
                       blu-scribbles                  [#051833 #0a4f66 #0f998e #12cc7f]
                       pen-n-paper                    [#e4dbba #a4929a #4f3a54 #260d1c]
                       2-bit-grayscale                [#000000 #676767 #b6b6b6 #ffffff]
                       nintendo-gameboy-bgb           [#081820 #346856 #88c070 #e0f8d0]
                       spacehaze                      [#f8e3c4 #cc3495 #6b1fb1 #0b0630]
                       pokemon-sgb                    [#181010 #84739c #f7b58c #ffefff]
                       nintendo-super-gameboy         [#331e50 #a63725 #d68e49 #f7e7c6]
                       rustic-gb                      [#2c2137 #764462 #edb4a1 #a96868]
                       ajstarks                       [#aa0000 #aaaaaa #000000 #ffffff]
                       blk-aqu4                       [#002b59 #005f8c #00b9be #9ff4e5]
                       dark-mode                      [#212121 #454545 #787878 #a8a5a5]
                       mist-gb                        [#2d1b00 #1e606e #5ab9a8 #c4f0c2]
                       red-brick                      [#eff9d6 #ba5044 #7a1c4b #1b0326]
                       moonlight-gb                   [#0f052d #203671 #36868f #5fc75d]
                       links-awakening-sgb            [#5a3921 #6b8c42 #7bc67b #ffffb5]
                       kankei4                        [#ffffff #f42e1f #2f256b #060608]
                       kirokaze-gameboy               [#332c50 #46878f #94e344 #e2f3e4]
                       nostalgia                      [#d0d058 #a0a840 #708028 #405010]
                       autumn-decay                   [#313638 #574729 #975330 #c57938 #ffad3b #ffe596]
                       polished-gold                  [#000000 #361c1b #754232 #cd894a #e6b983 #fff8bc #ffffff #2d2433 #4f4254 #b092a7]
                       funk-it-up                     [#e4ffff #e63410 #a23737 #ffec40 #81913b #26f675 #4c714e #40ebda #394e4e #0a0a0a]

```
