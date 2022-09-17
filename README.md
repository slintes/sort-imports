# sort-imports

Sort Go imports according to our best practices, separated by an empty line:

- Internal packages
- All other packages
- Openshift packages
- Kubernetes packages
- Packages of own module

## Usage

```
go install github.com/slintes/sort-imports
sort-imports <project_dir> [-w]
```

- the `-w` option will overwrite source files with fixed imports (use on your own risk,
use source control!) 
- the return code will indicate status:
  - `0`: all imports are fine or were overwritten
  - `1`: an error occured
  - `2`: there are unfixed unsorted or ungrouped imports, check output for details
