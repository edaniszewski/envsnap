
# envsnap

[![Build Status](https://travis-ci.org/edaniszewski/envsnap.svg?branch=master)](https://travis-ci.org/edaniszewski/envsnap)


`envsnap` is a small command-line utility for generating snapshots of a development
or runtime environment to aid in debugging and issue reporting.

It works by having projects define a [`.envsnap`](.envsnap) file, which defines various
data points that the project finds relevant. This ranges from OS/Arch to package dependencies.
`envsnap` generates a snapshot of the specified data points, making it easy for the user
to provide context for their environment when debugging or reporting an issue. 

With `envsnap`, the you no longer has to worry about which data points are relevant when
reporting an issue. Let the project define it, and let `envsnap` collect it.

## Installation

coming soon

## Usage

`envsnap` is a simple tool with only two commands:

* `envsnap init` - initializes a new `.envsnap` config
* `envsnap render` - render your environment based on the `.envsnap` config

For additional details and usage info, see the help info with `envsnap --help`.

### Example

```console
$ envsnap render
#### Environment

**System**
- _os_: darwin
- _arch_: x86_64
- _cpus_: 12
- _kernel_: Darwin
- _kernel version_: 19.0.0
- _processor_: i386

**Golang**
- _version_: go1.13.4
- _goroot_: /usr/local/Cellar/go/1.13.4/libexec
- _gopath_: /Users/edaniszewski/go


```

## Configuration

coming soon

## License

`envsnap` is released under the MIT license.
