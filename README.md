### Validation Error Messages

| Error | Text | Meaning |
|---|---------|----------|
| ErrTypeNotFound | type with name "{TypeName}" in "{ParentObject}" was not found | A type was referenced as value but not defined anywhere in the YAML document. |
| ErrInvalidValue | value assigned to key "{KeyName}" in "{ParentObject}" is invalid | An invalid value was defined (nil, "", List, Object in Object). |
| ErrIllegalTypeName | illegal type name "{KeyName}" in "{ParentObject}" | A type was named without adhering to go's syntax limitations (e.g. "fo$o", "func", "<-+"). |
| ErrRecursiveTypeUsage | illegal recursive type detected for "{RecurringKeyNames}" | A recursive type was defined. |
| ErrMissingPackageName | package name was not specified in the "_package" field at root level | A package name is required (e.g. "main"). |
| ErrInvalidPackageName | name "{PackageName}" is not a valid package name | An invalid name was assigned to the "_package" field. |
