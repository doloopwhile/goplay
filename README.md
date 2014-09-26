# gopark

Create local go sandbox.
You can edit source with favorite editor and install extra packages.

## Demo
Just run `gopark`.
Or, with option `-e` to open with editor.

## Install
```
go get github.com/doloopwhile/gopark
```

## Configuration

### `gopark.root` in `git config`
The path to directory in which `gopark` create sandbox directory.
Defaults to `/tmp/gopark`.

### `$EDITOR` environmnet variable
Editor to open created go file if `-e` option is specified.

## VS.
 - Go Playground(http://play.golang.org)

`gopark` create text environment in local machine.
Therefore, you can your favorite editor and install extra packages.

## License
It is MIT License. Please feel free to use, distribute or fork.
Please see LICENSE file for details.

## Contribution
I am looking for your pull request.
https://github.com/doloopwhile/gopark/pulls
