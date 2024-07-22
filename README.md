# zfmt

[![License](https://img.shields.io/github/license/zillow/zfmt)](https://github.com/zillow/zfmt/blob/main/LICENSE)
[![GitHub Actions](https://github.com/zillow/zfmt/actions/workflows/go.yml/badge.svg)](https://github.com/zillow/zfmt/actions/workflows/go.yml)
[![Codecov](https://codecov.io/gh/zillow/zfmt/branch/main/graph/badge.svg?token=STRT8T67YP)](https://codecov.io/gh/zillow/zfmt)

## Installation

`go get -u github.com/zillow/zfmt`

## About

A module which defines several concrete `Formatter` types responsible for serializing/deserializing objects.
The module centralizes this functionality and is leveraged by several zillow transport libs for use in configuration
driven serialization.
