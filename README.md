gopark
======

Create local go sandbox.
You can edit source with favorite editor and install extra packages.

Demo
====
Just run `gopark`.


Or, with option `-e` to open with editor.


Install
=======

```
go get github.com/doloopwhile/gopark
```

Configuration
=============

`gopark.root` in `git config`
-------------------------
The path to directory in which `gopark` create sandbox directory.
Defaults to `/tmp/gopark`.

`$EDITOR` environmnet variable
--------------------------
Editor to open created go file if `-e` option is specified.

VS.
===
 - Go Playground(http://play.golang.org)

`gopark` create text environment in local machine.
Therefore, allow you to use your favorite editor and install extra packages.
