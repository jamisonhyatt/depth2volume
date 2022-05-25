# depth2volume
dumb utility to calculate total volume to fill/excavate from a set of grade measurements.  

# How to measure
The measurements you would typically due with a laser level, and you would pass your "zero" measurement as `-initial-depth`.   Measurements greater than 
the zero value would necessitate fill being brought in to the area, while measurements less than the zero value would require removing soil. The measurements 
must be made in equidistant grids for this to work.  

# How does this work
It's dumb.  It operates on a grid, and for each square in the grid finds an average of all 4 measurements (1 at each corner) i.e., `(m1 + m2 + m3 + m4)/4` and then
subtracts the zero value to get a net volume needed for that square.

The total volume is just the sum of all the individual volumes from the squares

This, of course, means the more measurements you take, the more accurate this tool becomes.

Units don't matter, as long as they're all the same units.  

# Example usage
Calculate the volume needed to fill the [sample csv](./depths-sample.csv) with a 12ft x 27ft square (144 x 324 inches) at a starting depth of 1.25"

``` 
$> go run main.go \
  -iniital-depth 1.25 \
  -measurement-lineal 144 \
  -row-lineal 324 \
  -file depths-sample.csv
calculating cubic volume for 6 measurements
grid is 1 x 2 with a total area of 93312.000000 units
space between measurements is 144 units
space between rows is 324 units
initial depth is set to 1.25 units
Fill needed!
11080.800000 cubic units
```

Tells me I need 1108 cubic inches of fill brought in, which is almost 1/4 of a cubic yard (6.4 cubic feet) of fill for the area.
