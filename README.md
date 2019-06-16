<h2 align="center">Time to <i>Go Home</i> :clock10:</h2>
<p align="center">
	<a href="https://github.com/fgrosse/go-home/releases"><img src="https://img.shields.io/github/tag/fgrosse/go-home.svg?label=version&color=brightgreen"></a>
	<a href="https://github.com/fgrosse/go-home/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

---

Go Home is a small OpenGL based progress bar widget for your Desktop that
displays for how long you have been working each day. This is helpful for people
who tend to loose track of time and thus do overhours when they actually wanted
to leave home. I made this to learn a bit about simple OpenGL programming using
[the Go programming language][go].

<p align="center">
<img src="assets/screenshot_01.png">
<img src="assets/screenshot_02.png">
</p>

## Installation

### Precompiled binaries

You can find precompiled binaries at the [releases] page of the GitHub
repository.

### From Source

Go Home is packaged using [Go modules][go-modules]. Since this is not a library
but a runnable application within a `main` package you need to clone this
repository first. Typically this should be done outside of the `$GOPATH` or Go
will complain due Modules being enabled.

After you cloned the repo you should make sure to install the external
dependencies (i.e. OpenGL bindings) as explained at the
[GLFW repository][external-deps]. There is a [Makefile](Makefile) to install the
requires libraries on RedHead/Fedora.

Afterwards you simply use `go build` or `go install` and Go will fetch the
correct Go dependencies for you:
 
```bash
$ git clone https://github.com/fgrosse/go-home.git
Cloning into 'go-home'...
remote: Enumerating objects: 122, done.
remote: Counting objects: 100% (122/122), done.
remote: Compressing objects: 100% (69/69), done.
remote: Total 122 (delta 62), reused 109 (delta 49), pack-reused 0
Receiving objects: 100% (122/122), 71.62 KiB | 601.00 KiB/s, done.
Resolving deltas: 100% (62/62), done.

$ cd go-home               
$ make setup
…

$ go install && go-home --debug
2019-06-16 13:16	DEBUG	go-home/config.go:50	Running in debug mode
2019-06-16 13:16	INFO	go-home/config.go:54	Loading configuration	{"path": "/home/fgrosse/.go-home.yml"}
2019-06-16 13:16	INFO	go-home/app.go:56	Starting application	{"config": {"check_in": "2019-06-16 12:18", "work_duration": "8h0m0s", "lunch_duration": "1h0m0s", "day_end": "20:00"}}
```

## Usage

You can start the program without any arguments which will create an
undecorated window that displays when you started the program the first time
today and when its time to go home. The default configuration assumes you are
working 8 hours a day and do 1 hour of lunch break.

As time goes by the progress bar will slowly fill up from green to red. If you 
are working overtime it will start to pulse red to catch your attention. At this
point you should leave home and enjoy your free time with your family and
friends :relaxed:.

### Configuration

Go Home reads configuration from `$HOME/.go-home.yml`. If this file does not
exist on the first start it will be created using sensible default values.
The available options in there should be pretty self explanatory.

## Built With

* [pixel](https://github.com/faiface/pixel) - A hand-crafted 2D game library in Go
* [glfw](https://github.com/go-gl/glfw) - Go bindings for GLFW 3
* [cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions 
* [zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging in Go
* [pkg/errors](https://github.com/pkg/errors) - Simple error handling primitives
* [gopkg.in/yaml](https://gopkg.in/yaml.v3) - YAML support for the Go language 
* [Glacial Indifference Font](https://fontlibrary.org/en/font/glacial-indifference) - An open source typeface by Alfredo Marco Pradil ([SIL Open Font License](LICENSE_FONT))

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on the code of
conduct and on the process for submitting pull requests to this repository.

## Versioning

This software uses [SemVer] for versioning.
For the versions available, see the [tags on this repository][tags]. 

## Authors

- **Friedrich Große** - *Initial work* - [fgrosse]

See also the list of [contributors] who participated in this project.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

[releases]: https://github.com/fgrosse/go-home/releases
[external-deps]: https://github.com/go-gl/glfw/blob/master/README.md
[go]: https://golang.org
[go-modules]: https://github.com/golang/go/wiki/Modules
[SemVer]: http://semver.org
[tags]: https://github.com/fgrosse/go-home/tags
[fgrosse]: https://github.com/fgrosse
[contributors]: https://github.com/github.com/fgrosse/go-home/contributors
