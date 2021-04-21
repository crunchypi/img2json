# img2json
Simple CLI tool for converting images to JSON. Also has a couple filters: color range & random (didn't need anything else, sz).

# Examples
#### Original:
![alt text](https://raw.githubusercontent.com/crunchypi/img2coordinates/master/demo/cat.png?raw=true)

#### After filtering (pink is substitute for nothing (points are black)):
![alt text](https://github.com/crunchypi/img2coordinates/blob/master/demo/screenshot.png?raw=true)

#### Another example, animation in p5js https://editor.p5js.org/crunchypi/sketches/bGcjFyEEa:
![alt text](https://raw.githubusercontent.com/crunchypi/img2coordinates/master/demo/doge.gif)

# JSON format example.
  
    [
      {
          "x": 132,     # // x is x coordinate
          "y": 90,      # // y is y coordinate
          "r": 1,       # // r is r in rgba
          "g": 1,       # // g is g in rgba
          "b": 1        # // b is g in rgba
          "a": 255      # // a is a in rgba
       },
      ...,
    ]
      

# CLI args:
```
____________________________________________________________
Image to JSON converter with a bit of extra filtering stuff.
____________________________________________________________
 
This cli uses piping, meaning that argument positioning does
matter. The arguments are as follows (examples below):

    -help             Prints out this message.
    -ii      <path>   Specifies the 'Input Image' path.
    -ij      <path>   Specifies the 'Input JSON' path.
    -oi      <path>   Specifies the 'Output Image' path.
    -oj      <path>   Specifies the 'Output JSON' path.
    -byColor <...>    Accepts a value with the form:
                        int,int,int,int,int,int
                      .. which represents an RGB range.
                      So the first three ints are min RGB,
                      and the three latter are max RGB.
                      Everything outside of that range is
                      filtered out, i.e removed.
    -byRand  <int>    Filters a specified percentage (0..100)
                      of all currently stored pixels.
                  
Examples:
    .. To input an image, filter 50% of randomly selected
    pixels (uniform selection) and store as a new image:
        > -ii ./my.png -byRand 50 -oi ./my2.png

    .. To Input an image, filter out anything that is not
    absolute darkness, then store it as a JSON:
        > -ii ./my.png -byColor 0,0,0,0,0,0 -oj ./my2.json

    .. To import a JSON file created with this CLI and
    save it as an image:
        > -ij ./my.json -oi > -oi ./my.png
```
