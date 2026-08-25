# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0]

### Features

- write EPUB publications from Markdown and directories
- serialize publications to memory
- support `text/html` content and XML self-closing tags
- add `dcterms:modified`, infer SVG media type, and set spine toc
- support multiple package renditions
- expose landmarks and page-list navigation
- de-obfuscate embedded fonts

### Bug Fixes

- trim whitespace from mimetype file detection
- marshal `xml:lang` with the reserved XML namespace
- return an error when no cover image is defined
- fix nested navigation, duplicate ids, and title metadata
- capture metadata and rights inner XML content
- use the `xmldsig` namespace for signature elements
- propagate META-INF parse errors
- write mimetype first, uncompressed, and without an extra field
- normalize metadata values and resolve the unique identifier
- use correct XML namespace URIs in package and OCF documents

### Code Refactoring

- simplify root directory resolution

### Tests

- cover in-memory serialization
- cover guide reference type marshaling
- cover reader, metainf, and writer paths
- reorganize tests by source file and raise coverage to 90%
- port the W3C EPUB test suite with a nested structure
- add unique identifier resolution tests

### Documentation

- document the public API, add examples, and rewrite the README
- add CHANGELOG

## [0.3.0]

Initial tagged release.
