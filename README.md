# Tarsnap Prometheus Exporter

A [Tarsnap](https://www.tarsnap.com/) prometheus exporter that can be used

## Run

    $ export TARSNAP_EMAIL="xx@example.org"
    $ export TARSNAP_PASSWORD="password-here"
    $ go run main.go

## Docker

A docker image is published which can be used `drubin/tarsnap-prometheus-exporter:latest`

## Environment variables

The system depends on env variables for configuration

* `TARSNAP_EMAIL` - The Tarsnap account email
* `TARSNAP_PASSWORD` - The Tarsnap password account password

## Exported values

* `tarsnap_account_balance{account=""}` - Current account balance as latest data shows (generally 24 hours old)
