Package `generic` implements logic that should be applied for all
data types implemented as configuration parameter implementations.

The logic implemented here is not a good fit for, say, `Configurator` an interface
method because it either cannot or should not be applied at for any given
interface implementation. 

- [ ] Should this logic just move to package `convert`?

