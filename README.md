# krew-index-tracker
[![Build Status](https://travis-ci.com/corneliusweig/krew-index-tracker.svg?branch=master)](https://travis-ci.com/corneliusweig/krew-index-tracker)
[![LICENSE](https://img.shields.io/github/license/corneliusweig/krew-index-tracker.svg)](https://github.com/corneliusweig/krew-index-tracker/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/corneliusweig/krew-index-tracker)](https://goreportcard.com/report/corneliusweig/krew-index-tracker)
<!--
[![Code Coverage](https://codecov.io/gh/corneliusweig/krew-index-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/corneliusweig/krew-index-tracker)
[![Releases](https://img.shields.io/github/release-pre/corneliusweig/krew-index-tracker.svg)](https://github.com/corneliusweig/krew-index-tracker/releases)
-->

A tool to track download counts of `krew` plugins.

This tool does the following:

* It updates the copy of the `krew-index` (https://github.com/kubernetes-sigs/krew-index) in directory `index`.
* It reads all plugin manifests and looks for github URLs in the `spec.homepage` field.
* It fetches the release information for each repo via GitHub API.
* It writes the result into a BigQuery table.

