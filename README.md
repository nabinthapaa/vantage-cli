# Introduction

This is a simple project to mimic some behaviour of Lenovo vantage from windows.

## Features

- [x] Turn on or off Conservation mode
- [x] Turn on or off Function key lock
- [x] Turn on or off USB Charging

## Commands

###### ! You need to run this command as super user

###### ! Access file from /sys/bus/platform/devices/VPC2004:00/\*

```bash
    vantage on # Turns on Conservation mode
    vantage off # Turns off Conservation mode
    vantage table # Shows a table layout to turn on off using enter key
```

## Installation

Clone this repo and run make install

```bash
  git clone https://www.github.com/nabinthapaa/vantage-cli
  sudo make install # for moving the file to the /usr/bin/ directory
  # Run above commands
```
