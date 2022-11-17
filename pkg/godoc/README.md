Example showcasing Go documentation with `godoc`.
Read more in "Efficient Go"; Example 2-9.

## Usage

Prerequisite: `go` installed.

1. Pull the repo on local machine.
2. Go to the directory where this README exists.
3. Run `make`. This will trigger `godoc` installation on your local machine. Then it will start `godoc` in serving mode that serves documentation on port 6060.
4. Once you see "using module mode; ...", go to `localhost:6060` in your browser. You see all available modules on your machine, including standard library. For `pkg/godoc` package see [http://localhost:6060/pkg/github.com/efficientgo/examples/pkg/godoc/](http://localhost:6060/pkg/github.com/efficientgo/examples/pkg/godoc/).
