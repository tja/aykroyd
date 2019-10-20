<p align="center"><img width="200" src="web/images/hero.svg"></a></p>

<p align="center">
  <a href="https://godoc.org/github.com/tja/aykroyd"><img src="https://godoc.org/github.com/tja/aykroyd?status.svg" alt="GoDoc"></a>
  <a href="https://goreportcard.com/report/github.com/tja/aykroyd"><img src="https://goreportcard.com/badge/github.com/tja/aykroyd" alt="Go Report Card"></a>
  <a href="http://opensource.org/licenses/MIT"><img src="http://img.shields.io/badge/license-MIT-brightgreen.svg" alt="MIT License"></a>
</p>


# Aykroyd &mdash; Email forwards via Postfix

**Aykroyd** is an HTTP server and web application that allows the user to manage email forwards in
[Postfix](http://www.postfix.org).


## Installation

### Binaries

Pre-built binaries are available on the [release page](https://github.com/tja/aykroyd/releases/latest). Simply
download, make executable, and move it to a folder in your `PATH`:

```bash
curl -L https://github.com/tja/aykroyd/releases/download/v2.0.0/aykroyd-`uname -s`-`uname -m` >/tmp/aykroyd
chmod +x /tmp/aykroyd
sudo mv /tmp/aykroyd /usr/local/bin/aykroyd
```

### Docker

**Aykroyd** is also available as a [Docker image from Docker Hub](https://hub.docker.com/r/neathack/aykroyd).
Configuration can be done via a config file or environment variables (see below for details).


## Setup

### MariaDB / MySQL

Install [MariaDB](https://mariadb.com/downloads/) or [MySQL](https://dev.mysql.com/downloads/), create a
schema `postfix`, and grant a user access to it:

```mysql
CREATE SCHEMA `postfix`;
GRANT SELECT ON `postfix`.* TO `postfix`@`localhost` IDENTIFIED BY '<password>';
```

### Postfix

Install [Postfix](http://www.postfix.org) with [MySQL support](http://www.postfix.org/MYSQL_README.html). Here
is an example for Debian/Ubuntu:

```bash
sudo apt -y install postfix postfix-mysql
```

Add the Postfix domain config file `/etc/postfix/mysql_virtual_alias_domains.cf`:

```
user = postfix
password = <password>
hosts = 127.0.0.1
dbname = postfix
table = domains
select_field = name
where_field = name
```

Add the Postfix forwards config file `/etc/postfix/mysql_virtual_alias_forwards.cf`:

```
user = postfix
password = <password>
hosts = 127.0.0.1
dbname = postfix
table = forwards
select_field = to
where_field = from
```

Append the references to the two aforementioned files to `/etc/postfix/main.cf`:

```
virtual_alias_domains = mysql:/etc/postfix/mysql_virtual_alias_domains.cf
virtual_alias_maps = mysql:/etc/postfix/mysql_virtual_alias_forwards.cf
```


## Usage

### Binary

Run `aykroyd` in the command line like this:

```bash
aykroyd --host "0.0.0.0:8080"
```

Open a browser and visit `http://localhost:8080` to bring up the web interface. See below for how to configure
host, port, and database access.

### Docker

Run the Docker image `neathack/aykroyd` and proxy port `80` to a port of your choice:

```bash
docker run -ti --rm -p "8080:80" neathack/aykroyd:latest
```

Open a browser and visit `http://localhost:8080` to bring up the web interface. See below for how to configure
database access.


## Configuration

**Aykroyd** can be configured via command line switches, a configuration file (in the current folder, in
`~/.config/aykroyd/config.yml`, or in `/etc/aykroyd/config.yml`), or environment variables. The following
options are available:

| Command Line          | Config File   | ENV Variable          | Description                                                  |
|-----------------------|---------------|-----------------------|--------------------------------------------------------------|
| `-l`, `--listen`      | `listen`      | `AYKROYD_LISTEN`      | IP and port the server will listen on (default "0.0.0.0:80") |
| `-a`, `--asset`       | `asset`       | `AYKROYD_ASSET`       | Path to static web assets (default uses embedded assets)     |
| `-H`, `--db-host`     | `db-host`     | `AYKROYD_DB_HOST`     | MySQL host (default "localhost")                             |
| `-d`, `--db-database` | `db-database` | `AYKROYD_DB_DATABASE` | MySQL database (default "postfix")                           |
| `-u`, `--db-username` | `db-username` | `AYKROYD_DB_USERNAME` | MySQL username (default "postfix")                           |
| `-p`, `--db-password` | `db-password` | `AYKROYD_DB_PASSWORD` | MySQL password                                               |
| `-v`, `--verbose`     | `verbose`     | `AYKROYD_DB_VERBOSE`  | Write more                                                   |


## License

Copyright (c) 2018&ndash;2019 Thomas Jansen. Released under the
[MIT License](https://github.com/tja/aykroyd/blob/master/LICENSE).

Email icon made by [Pixel Buddha](https://www.flaticon.com/authors/pixel-buddha) from
[www.flaticon.com](https://www.flaticon.com/) is licensed by
[CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/).
