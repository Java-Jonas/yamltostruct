## Validation Error Messages
<br/> 

### structural:

| Error | Text | Meaning |
|---|---------|----------|
| ErrIllegalValue | value assigned to key "{KeyName}" in "{ParentObject}" is invalid | An invalid value was defined (nil, "", List, Object in Object). |
| ErrMissingPackageName | package name was not specified in the "_package" field at root level | A package name is required (e.g. "main"). |
<br/> 

### syntactical:
| Error | Text | Meaning |
|---|---------|----------|
| ErrIllegalTypeName | illegal type name "{KeyName}" in "{ParentObject}" | A type was named without adhering to go's syntax limitations (e.g. "fo$o", "func", "<-+"). |
| ErrIllegalPackageName | name "{PackageName}" is not a valid package name | An invalid name was assigned to the "_package" field. |
| ErrInvalidValueString | value "{ValueString}" assigned to "{KeyName}" in "{ParantObject}" is invalid | An invalid value was assigned to a key |
<br/> 

### logical:
| Error | Text | Meaning |
|---|---------|----------|
| ErrTypeNotFound | type with name "{TypeName}" in "{ParentObject}" was not found | A type was referenced as value but not defined anywhere in the YAML document. |
| ErrRecursiveTypeUsage | illegal recursive type detected for "{RecurringKeyNames}" | A recursive type was defined. |
| ErrInvalidMapKey | "{MapKey}" is not a valid map key | An uncomparable type was chosen as map key. |
<br/> 

### TODO:
- recursiveness check for map and slice values
- allow only use of comparable types as map keys
- package name validation
