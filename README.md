Strip JSON Comments
-------------------
This is a Go module version of the go script from the original [repo](https://github.com/dnmfarrell/Strip-JSON-Comments).

Strips single line `//` and block `/* */` comments from .jsonc text:

```
echo '[1, /* strip this comment */ 2, 3]' | ./stripjsonc
[1,                          2, 3]
```

### Comments are replaced with spaces
To preserve line and column locations of the remaining JSON text.

### Block comments do not nest
JSONC doesn't have a spec but its block comments have been [described](https://code.visualstudio.com/docs/languages/json#_json-with-comments) as behaving like JavaScript, i.e. they do not nest.

### Trailing commas are not removed
Some JSONC parsers support trailing commas, but as `stripjsonc` is line-oriented, it does not have an option to remove trailing commas.

### Test input
This repo contains a test file called `input.jsonc` which exercises the different cases of embedding comments in JSON. The file `output.jsonc` contains the stripped output.

```
./stripjsonc < input.jsonc
                                
{ "foo //": 123,                               
     
                  
        
   
     
                                 
     "b😃\"\\ar/*":                                   "*/baz"
}
```

### Building
To create the `stripjsonc` program, simply run:

   go build

### Fuzzing
This module comes with a fuzzing test, which can be run with `go test`:

    go test -fuzz=Fuzz -fuzztime=30s
