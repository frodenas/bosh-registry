# BOSH Registry

This is a **simple** and **experimental** [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client and Server.

## Disclaimer

This is **NOT** presently a production ready BOSH Registry. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## BOSH Registry Client

### Installation

Using the standard `go get`:

```
$ go get github.com/frodenas/bosh-registry
```

### Usage

For usage and examples see the [Godoc](https://godoc.org/github.com/frodenas/bosh-registry/client).

## BOSH Registry Server

### Installation

Using the standard `go get`:

```
$ go get github.com/frodenas/bosh-registry
```


### Usage

Create a configuration file:

```JSON
{
  "server": {
    "protocol": "http",
    "address": "127.0.0.1",
    "port": 25777,
      "username": "admin",
      "password": "admin",
    "tls": {
      "certfile": "certfile.pem",
      "keyfile": "keyfile.pem",
      "cacertfile": "ca.crt"
    }
  },
  "store": {
    "adapter": "bolt",
    "options": {
      "dbfile": "registry.db"
    }
  }
}
```

Run the registry using the previously created configuration file:

```
$ registry -configPath="Path to configuration file"
```

## Contributing

In the spirit of [free software](http://www.fsf.org/licensing/essays/free-sw.html), **everyone** is encouraged to help improve this project.

Here are some ways *you* can contribute:

* by using alpha, beta, and prerelease versions
* by reporting bugs
* by suggesting new features
* by writing or editing documentation
* by writing specifications
* by writing code (**no patch is too small**: fix typos, add comments, clean up inconsistent whitespace)
* by refactoring code
* by closing [issues](https://github.com/frodenas/bosh-registry/issues)
* by reviewing patches

### Submitting an Issue
We use the [GitHub issue tracker](https://github.com/frodenas/bosh-registry/issues) to track bugs and features.
Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate
support for an existing issue by voting it up. When submitting a bug report, please include a
[Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug,
including your gem version, Ruby version, and operating system. Ideally, a bug report should include a pull request with
 failing specs.

### Submitting a Pull Request

1. Fork the project.
2. Create a topic branch.
3. Implement your feature or bug fix.
4. Commit and push your changes.
5. Submit a pull request.

## Copyright

Copyright (c) 2015 Ferran Rodenas. See [LICENSE](https://github.com/frodenas/bosh-registry/blob/master/LICENSE) for details.
