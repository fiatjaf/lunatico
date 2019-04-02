This library makes it easy for people who don't understand the Lua stack thing to run Lua scripts inside Go programs. It uses https://github.com/aarzilli/golua underneath (but could be easily ported to other Lua implementations as the functions are almost the same).

See [the small test file](lunatico_test.go) for how to use (hint: it's basically (i) set some globals, (ii) run Lua code, (iii) get globals back).
