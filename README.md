
# FuguFugu

![FuguFugu](https://github.com/inkmi/fugufugu/blob/main/Logo.png?raw=true)

!! EXPERIMENTAL !!

Created for safe migrations of [Inkmi - Dream Jobs for CTOs](https://www.inkmi.com)

Detected destructive or incompatible changes:
* `DROP TABLE table`
* `DROP VIEW view`
* `ALTER TABLE table DROP COLUMN column`

## Installation

```bash
go install github.com/inkmi/fugufugu@latest
```

## Run

Run `fugufugu` on the migrations directory to check your SQL migrations

```bash
> fugufugu --dir migrations/
```

