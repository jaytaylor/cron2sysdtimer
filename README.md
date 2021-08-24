# jaytaylor.com/cron2sysdtimer

Automatically transform crontab jobs into systemd timers

[![Documentation](https://godoc.org/github.com/jaytaylor/cron2sysdtimer?status.svg)](https://godoc.org/github.com/jaytaylor/cron2sysdtimer)
[![Build Status](https://travis-ci.org/jaytaylor/cron2sysdtimer.svg?branch=master)](https://travis-ci.org/jaytaylor/archiveis)
[![Report Card](https://goreportcard.com/badge/github.com/jaytaylor/cron2sysdtimer)](https://goreportcard.com/report/github.com/jaytaylor/cron2sysdtimer)

## Requirements

* Golang v1.17 or newer
* Linux with SystemD (naturally :)

## Installation

```bash
go install jaytaylor/cron2sysdtimer@latest
```

## Usage

cron2sysdtimer reads crontab file at `/etc/crontab` by default. You can specify crontab file with `-f FILE` flag.

systemd unit file are saved at `/run/systemd/system` by default. You can specify save directory with `-o OUTDIR` flag.

```bash
$ cron2sysdtimer
$ cron2sysdtimer -f sample.cron -o unitfiles
```

### Example

```bash
ubuntu@ubuntu-jovial:~/src/jaytaylor.com/cron2sysdtimer$ sudo cron2sysdtimer -f sample.cron --reload
ubuntu@ubuntu-jovial:~/src/jaytaylor.com/cron2sysdtimer$ systemctl list-timers
NEXT                         LEFT                   LAST PASSED UNIT                         ACTIVATES
Fri 2017-01-20 07:50:00 UTC  4min 16s left          n/a  n/a    cron-77e2fb273c45.timer      cron-77e2fb273c45.service
Fri 2017-01-20 07:56:01 UTC  10min left             n/a  n/a    systemd-tmpfiles-clean.timer systemd-tmpfiles-clean.service
Fri 2017-01-20 08:00:00 UTC  14min left             n/a  n/a    cron-1b33d99b7dda.timer      cron-1b33d99b7dda.service
Fri 2017-01-20 10:00:00 UTC  2h 14min left          n/a  n/a    cron-b60fe106ef63.timer      cron-b60fe106ef63.service
Fri 2017-01-20 12:16:09 UTC  4h 30min left          n/a  n/a    snapd.refresh.timer          snapd.refresh.service
Fri 2017-01-20 19:11:59 UTC  11h left               n/a  n/a    apt-daily.timer              apt-daily.service
Wed 2017-02-01 00:00:00 UTC  1 weeks 4 days left    n/a  n/a    cron-fcd6d8377d9d.timer      cron-fcd6d8377d9d.service
Sat 2017-12-02 01:23:00 UTC  10 months 11 days left n/a  n/a    cron-d3c507cb2439.timer      cron-d3c507cb2439.service

8 timers listed.
Pass --all to see loaded but inactive timers, too.
```

### Reload systemd and start all timers automatically

If `--reload` is provided, cron2sysdtimer reloads systemd unit files (= `systemctl daemon-reload`) and starts all generated timers (= `systemctl start foo.timer`). Maybe `sudo` is required to execute.

```bash
$ sudo cron2sysdtimer -f sample.cron --reload
```

### Determine unit name from command to execute

Crontab does not have the concept of "task name", and a task name is required to identify each systemd unit.

`cron2sysdtimer` supports automatically generating the task name from the original command using a regular expression- Use the `--name-regexp REGEXP` flag.

The regex must have one [capturing group](http://www.regular-expressions.info/brackets.html).

If a regular expression is not provided or command does not match to the given regular expression, a hash value is automatically calculated and used as the unit name.

```bash
$ cron2sysdtimer -f sample.cron --name-regexp '--name ([a-zA-Z0-9_-]+)'
```

### Delete unregistered unit files

If `--delete` is provided, cron2sysdtimer deletes unit files which are no longer written in the given crontab file.

```bash
$ cron2sysdtimer -f tmp/scheduler -o /run/systemd/system --delete
Deleted: /run/systemd/system/cron-19fb9c164fe8.service
Deleted: /run/systemd/system/cron-19fb9c164fe8.timer
Deleted: /run/systemd/system/cron-4f76a3902132.service
Deleted: /run/systemd/system/cron-4f76a3902132.timer
```

### Specify unit dependencies

You can specify unit dependencies (`After=`) with `--after AFTER` flag.

```bash
$ cron2sysdtimer -f sample.crom --after docker.service
```

## Development

### Get project build dependencies

```bash
go get -u github.com/jteeuwen/go-bindata/...
```

### Run tests

Standard:

```bash
cd $GOPATH/src/jaytaylor.com/cron2sysdtimer
go test -v ./...; echo $?
```

### Generate static bindata assets

```bash
cd $GOPATH/src/jaytaylor.com/cron2sysdtimer
go generate ./...
```

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Credits

### Original Author

Forked from [@dtan4's](https://github.com/dtan4) [ct2stimer.git](https://github.com/dtan4/ct2stimer).

