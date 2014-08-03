fractalserver
=============

A Google AppEngine program that serves images of the Mandelbrot set. Written in Go.

Responds to `GET /mandelbrot/` requests with a 256x256 pixel png image showing a region of the Mandlbrot set. Query parameters are:

* `x, y` the coordinate of the center pixel. (x + yi).
* `zoom` the range of the image. Distance from the center to the edge is 1 / zoom.

See it in action at [mandelbrotserver.appspot.com](http://mandelbrotserver.appspot.com).
