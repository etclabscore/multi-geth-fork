This package and included logic should be removed eventually.

It represents vestigal constructions for converting between chain configuration
data types. I believe all functions are used only in the cmd/puppeth package.

These functions should be replaced in their occurences with `convert.Convert` logic instead,
and then this package can be removed entirely.