# crta-box

[![SLSA Go releaser](https://github.com/omega1x/crta-box/actions/workflows/go-ossf-slsa3-publish.yml/badge.svg)](https://github.com/omega1x/crta-box/actions/workflows/go-ossf-slsa3-publish.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Industrial monitoring systems for power plants. Stream data from acoustic-based *culvert rupture telltale aggregation* boxes (*CRTA-BOX*es) (to target [ClickHouse](https://clickhouse.com/)-database).

## Usage

### Basics

To start streaming data from *CRTA-BOX* to *ClickHouse*-database execute

```shell
./crta-box stream <access-options>
```

where a full set of `<access-options>` could be listed by executing

```shell
./crta-box stream --help
```

Optionally they could check access to *CRTA-BOX* with

```shell
./crta-box box <access-options>
```

or *ClickHouse*-database with

```shell
./crta-box house <access-options>
```

where appropriate `<access-options>` could be listed with `--help`:

```shell
./crta-box box --help && ./crta-box house --help
```

### Logging

Enforce logging to file by adding `--log=<FILE>` option before command:

```bash
./crta-box --log=crta-box.log stream <access-options>
```

## Installation

### Prerequisites

For operability of `crta-box`, it is necessary not only to have valid access options but also the correct organization of the table structure in both communicating systems. The structure of the tables in the *CRTA-BOX* is determined by the current version of the installed telltale boxes. While the data structure in *ClickHouse* can be drawn up from a [ch__create-table__log_box3.sql](share/ch__create-table__log_box3.sql). In order to provide erroneous execution of `crta-box` there must be some data in *ClickHouse*-database that could be inserted by provided [ch__insert-table__log_box3.sql](share/ch__insert-table__log_box3.sql).

### Installation process

Download the latest binary from [Release page](https://github.com/omega1x/crta-box/releases/new):

```shell
wget <link> -o crta-box && chmod +x crta-box
```

Then check installation:

```shell
./crta-box --version
```
