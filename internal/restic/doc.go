// Package restic implements calling the restic executable.
//
// As far as possible restic is configured using the environment variables
// defined in restic's documentation
// (https://restic.readthedocs.io/en/stable/040_backup.html#environment-variables).
// This makes it easy to support a wide range of the features provided by
// resitc without writing too much custom code. On the downside this leads to
// a slightly peculiar for this package, as some functions require certain
// environment variables to be set to work correctly.
package restic
