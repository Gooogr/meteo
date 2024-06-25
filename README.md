# meteo
A CLI app for weather prediction in Go

## Installation

Create `config.yaml` based on `config\config.yaml.example` file.

Set path to Meteo config folder
`export METEO_CONFIG_PATH="/home/grigoriy/repos/meteo/config/"`

Run `go install`

## Usage
Get predictions in default coordinates
`meteo`

Example of getting weather on Null Island
`meteo --lat 0.0 --lon 0.0`