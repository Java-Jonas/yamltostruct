### Validation Error Messages

| Error | Text | Meaning |
|---|---------|----------|
| ErrTypeNotFound | type with name "{TypeName}" in "{ParentObject}" was not found | A type was referenced as value but not defined anywhere in the YAML document. |
| ErrInvalidValue | value assigned to key "{KeyName}" in "{ParentObject}" is invalid | An invalid value was defined (nil, "", List, Object in Object). |
| ErrIllegalTypeName | illegal type name "{KeyName}" in "{ParentObject}" | A type was named without adhering to go's syntax limitations (e.g. "fo$o", "func", "<-+"). |
| ErrRecursiveTypeUsage | illegal recursive type detected for "{RecurringKeyNames}" | A recursive type was defined. |
| ErrMissingPackageName | package name was not specified in the "_package" field at root level | A package name is required (e.g. "main"). |
| ErrInvalidPackageName | name "{PackageName}" is not a valid package name | An invalid name was assigned to the "_package" field. |

### TODO:
- field type syntax validation
- better naming for file and validation methods
- recursiveness check for map and slice values
- allow only use of comparable types as map keys
- package name validation
- order validation in:
    1. structural: no lists, objects in objects etc.
    2. syntactical: no illegal type names, values
    3. logical: no undefined types, recursiveness, non-comparable map keys etc.
