[![Build Status](https://github.com/google/renameio/workflows/Test/badge.svg)](https://github.com/google/renameio/actions?query=workflow%3ATest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/google/renameio)](https://pkg.go.dev/github.com/google/renameio)
[![Go Report Card](https://goreportcard.com/badge/github.com/google/renameio)](https://goreportcard.com/report/github.com/google/renameio)

The `renameio` Go package provides a way to atomically create or replace a file or
symbolic link.

## Atomicity vs durability

`renameio` concerns itself *only* with atomicity, i.e. making sure applications
never see unexpected file content (a half-written file, or a 0-byte file).

As a practical example, consider https://manpages.debian.org/: if there is a
power outage while the site is updating, we are okay with losing the manpages
which were being rendered at the time of the power outage. They will be added in
a later run of the software. We are not okay with having a manpage replaced by a
0-byte file under any circumstances, though.

## Advantages of this package

There are other packages for atomically replacing files, and sometimes ad-hoc
implementations can be found in programs.

A naive approach to the problem is to create a temporary file followed by a call
to `os.Rename()`. However, there are a number of subtleties which make the
correct sequence of operations hard to identify:

* The temporary file should be removed when an error occurs, but a remove must
  not be attempted if the rename succeeded, as a new file might have been
  created with the same name. This renders a throwaway `defer
  os.Remove(t.Name())` insufficient; state must be kept.

* The temporary file must be created on the same file system (same mount point)
  for the rename to work, but the TMPDIR environment variable should still be
  respected, e.g. to direct temporary files into a separate directory outside of
  the webserver’s document root but on the same file system.

* On POSIX operating systems, the
  [`fsync`](https://manpages.debian.org/stretch/manpages-dev/fsync.2) system
  call must be used to ensure that the `os.Rename()` call will not result in a
  0-length file.

This package attempts to get all of these details right, provides an intuitive,
yet flexible API and caters to use-cases where high performance is required.

## Major changes in v2

With major version renameio/v2, `renameio.WriteFile` changes the way that
permissions are handled. Before version 2, files were created with the
permissions passed to the function, ignoring the
[umask](https://en.wikipedia.org/wiki/Umask). From version 2 onwards, these
permissions are further modified by process' umask (usually the user's
preferred umask).

If you were relying on the umask being ignored, add the
`renameio.IgnoreUmask()` option to your `renameio.WriteFile` calls when
upgrading to v2.

## Windows support

It is [not possible to reliably write files atomically on
Windows](https://github.com/golang/go/issues/22397#issuecomment-498856679), and
[`chmod` is not reliably supported by the Go standard library on
Windows](https://github.com/google/renameio/issues/17).

As it is not possible to provide a correct implementation, this package does not
export any functions on Windows.

## Disclaimer

This is not an official Google product (experimental or otherwise), it
is just code that happens to be owned by Google.

This project is not affiliated with the Go project.
