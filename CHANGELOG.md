# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.3] - 2019-06-09

### Fixed
- Fix the `CMD` of the Docker image

## [0.2.2] - 2019-06-09

### Changed
- Add a multi-stage Docker build
- Add a `go get` step to the Docker build to leverage caching

## [0.2.1] - 2019-06-08

### Removed
- No longer update the deployment from `kube/deployment.yaml` since CI updates the deployment in k8s

## [0.2.0] - 2019-06-08

### Added
- Continuous deployment
- Continuous testing
