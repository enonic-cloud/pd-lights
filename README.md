<h1>pd-lights</h1>

[![build](https://github.com/enonic-cloud/pd-lights/actions/workflows/build.yml/badge.svg)](https://github.com/enonic-cloud/ec-backup-agent/actions/workflows/build.yml)

A small program that polls pagerduty and controls the traffic lights in the office

- [Running](#running)
- [Configuration file](#configuration-file)
- [Release](#release)

## Running

```console
$ pd-lights -h
A small app that controls the office traffic lights

Usage:
  pd-lights [flags]

Flags:
  -h, --help               help for pd-lights
      --ip string          npc ip
      --loop duration      loop interval (default 30s)
      --timeout duration   request timeout (default 30s)
      --token string       pagerduty token
```

## Configuration file

The program will look for a file called `.pd-lights.yaml` in the working directory and the `$HOME` directory. The file
can look something like this:

```yaml
token: ??
ip: ?.?.?.?
loop: 30s
timeout: 15s
```

## Release

The release procedure is controlled by github actions. Just tag a commit and a release will automatically be created for
you:

```console
$ git tag -a vX.Y.Z -m vX.Y.Z
$ git push origin vX.Y.Z
```
