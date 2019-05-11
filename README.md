# multiwatch

[![Build Status](https://travis-ci.org/Enapiuz/multiwatch.svg?branch=master)](https://travis-ci.org/Enapiuz/multiwatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/Enapiuz/multiwatch)](https://goreportcard.com/report/github.com/Enapiuz/multiwatch)
[![codecov](https://codecov.io/gh/Enapiuz/multiwatch/branch/master/graph/badge.svg)](https://codecov.io/gh/Enapiuz/multiwatch)
[![Maintainability](https://api.codeclimate.com/v1/badges/61bf67df2cdf15e5262f/maintainability)](https://codeclimate.com/github/Enapiuz/multiwatch/maintainability)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


Simple task runner on directory changes.

## Installation
### Manual
1. Download multiwatch
    - `git clone https://github.com/Enapiuz/multiwatch.git`
2. Install via go
    - `cd multiwatch && go install`
    
### Distros
Work In Progress

## Config
```toml
# debounce time for file change events
delay=500

[[watch]]
name = "linter"
paths = ["src"]
commands = ["npm run lint"]

[[watch]]
name = "tests"
paths = ["src", "tests"]
commands = ["npm run test"]
```
