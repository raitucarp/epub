# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Version
headings follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Features

- **Markdown authoring.** Add `Writer.AddMarkdown`, `Writer.AddMarkdownFile`,
  and `Writer.AddMarkdownDirectory` to convert Markdown sources into XHTML
  content documents. `AddMarkdownDirectory` reads every Markdown file in a
  directory, orders them by file name, and derives a nested table of contents
  from the `h1` through `h6` headings in each file.
- **In-memory serialization.** Add `Writer.WriteBytes` and `Writer.Build` for
  producing an EPUB without touching the filesystem, and `OCFZipContainer.WriteArchive`
  for writing the container to an arbitrary `io.Writer`.
- **Reader accessors.** Add `Reader.ReadContentMarkdownByHref`,
  `Reader.ReadImageBytesById`, and `Reader.ReadImageBytesByHref`.
- **EPUB 2 content compatibility.** Recognize `text/html` content documents and
  expand XML-style self-closing tags before parsing, so publications produced by
  EPUB 2 tooling parse correctly.
- **Multiple renditions.** Support publications with more than one package
  document via `Reader.ListRenditions`, `Reader.SelectPackageRendition`, and
  `Reader.CurrentSelectedPackage`.
- **Navigation.** Expose the landmarks and page-list navigation through
  `Reader.Landmarks` and `Reader.PageList`.
- **Font de-obfuscation.** Restore fonts obfuscated with the IDPF and Adobe
  algorithms described by the OCF specification.
- **Metadata.** Set the required `dcterms:modified` property automatically and
  expose a `Writer.Modified` method for explicit control.
- **Documentation.** Add Godoc examples that progress from a basic round-trip to
  Markdown-driven builds, and runnable programs under [`./examples`](./examples).

### Bug Fixes

- Accept a `mimetype` file with trailing whitespace instead of rejecting the
  container as a media-type mismatch.
- Marshal `xml:lang` using the reserved XML namespace instead of emitting an
  invalid `_xml` prefix.
- Resolve image hrefs correctly for resources nested more than one directory
  deep.
- Return an error from `Reader.CoverBytes` when the publication has no cover
  instead of panicking.
- Build nested NCX navigation correctly instead of flattening it or panicking
  on deep headings.
- Ensure manifest item IDs are unique and that a generated navigation document
  does not overwrite an existing content document with the same file name.
- Write additional titles as `dc:title` elements rather than malformed `meta`
  refinements.
- Propagate META-INF parse errors instead of silently discarding them.
- Parse `signatures.xml` elements with the `xmldsig` namespace.
- Capture the raw inner XML of `metadata.xml` and `rights.xml`.
- Write the `mimetype` file first, uncompressed, and without an extra field.
- Normalize metadata values and resolve the unique identifier from the package
  `unique-identifier` attribute.
- Use the correct XML namespace URIs throughout the package and OCF documents.
- Stop requiring a cover image; it is optional in EPUB 3.3.

### Code Refactoring

- Simplify `OCFZipContainer` root-directory resolution to use the `path` package
  and remove platform-specific dead code.
- Consolidate Markdown conversion into a single helper shared by the per-document
  and per-directory APIs.

### Tests

- Port the W3C EPUB 3.3 test suite with a nested, per-category structure.
- Reorganize tests to mirror the source layout and raise statement coverage to
  at least 90% for the root package and 95% for the `ocf`, `pkg`, and `ncx`
  packages.
- Add round-trip tests that read and rewrite a corpus of real-world publications.

## [0.3.0]

Initial tagged release.
