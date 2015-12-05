# caddyext

*caddyext* is a command line tool to manage *extensions/directives* of your [Caddy](http://caddyserver.com) build.

### Requirements

- Go (v1.4 or higher)
- [Caddy](http://github.com/mholt/caddy)'s source inside your GOPATH.

## Installing

    $ go get -u github.com/caddyserver/caddyext

## Usage

### Installing extension from Caddy's registry

[Caddy's registry](https://github.com/caddyserver/buildsrv/blob/master/features/registry.go)

Example:

    $ caddyext install git

### Installing extension from repository

Example:

    $ caddyext install search github.com/pedronasser/caddy-search

### Removing extension

Example:

    $ caddyext remove search

### Disabling core directive

Example:

    $ caddyext disable browse

### Moving extension on Caddy's stack

    $ caddyext stack

    Available Caddy directives/extensions:
       (✓) ENABLED | (-) DISABLED

       1. (✓) root (core)
       2. (✓) tls (core)
       ...
       21. (-) search

    $ caddyext move search 2

### Using different target caddy source

    $ CADDYPATH={import_path_to_caddy} caddyext ...

### More

```
Caddy's directive/extension manager

Usage:
  caddyext [command]

Available Commands:
  build       Build caddy from the current state
  install     Install and enables a extension
  remove      Remove an extension from caddy's directives source (only 3rd-party)
  stack       Show stack of directives/extensions
  enable      Enables a installed directive or extension
  disable     Disables a installed directive or extension
  move        Move target's index on Caddy's stack
  version     Show caddyext's version
  help        Help about any command


Global Flags:
  -h, --help=false: help for caddyext

Use "caddyext help [command]" for more information about a command.
```

## License

```
The MIT License (MIT)

Copyright (c) 2015 Pedro Nasser

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
