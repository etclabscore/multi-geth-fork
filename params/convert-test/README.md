Package `convert_test` is its own package because adequate testing requires the import of
several packages, `parity`, `goetherem`, and `multigeth`. 

The actual `convert` package should be agnostic of these data types and logics implementations;
only the testing should care about these specifics.
