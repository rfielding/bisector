A simple bisection library and command line
===========================================

I frequently wish that I had a Unix bisect utility that doesn't have anything to do with git.
Sometimes I find myself debugging something where I know some endpoints:

> f(203)  pass
> f(1033) fail

So, I want to fire up an interactive command line utility and note that 203 passes and 1033 fails:

> n 203 p 1033 f

It can then respond with a number in the middle to try, while still noting what the original range was:

> 503 203 p 1033 f

I then either type a 'p' or an 'f' in response to each line until it notes that it tried two numbers next to each other, and they marked the change from p to f:

> a 535 p 536 f

In most cases, I am checking from pass to fail, and want to know what introduced the fail, so in this case I choose the second number - 536.
This handles the more general case where I can note that it used to fail, and now it passes:

> f(600) fail
> f(993) pass

And just to be more general, it is handling the case where I picked a higher number first:

> f(999) pass
> f(603) fail

The requirement for the answer is that it returns a pair of adjacent numbers where it changed from pass to fail.

