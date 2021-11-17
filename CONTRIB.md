Contributing and editing some of the code requires some additional steps beyond simply building LaForge.

You'll need to install three tools:
* EasyJSON
* Renum (0.0.6)
* Fileb0x

```
go get -u github.com/mailru/easyjson/...
go get github.com/gen0cide/renum/cmd/renum
go get -u github.com/UnnoTed/fileb0x
```

If you make changes to a builder (TFGCP), you will need to run `go generate` to repackage the template files using filebox.
If you make changes within the core directory, be sure to run `go generate` to rebuild some JSON mappings and generated enums.

If you get the following error when running LaForge after a build:

`panic: gob: registering duplicate types for "github.com/zclconf/go-cty/cty.primitiveType": cty.primitiveType != cty.primitiveType`

The following steps can be taken to fix it:

1. Edit `$GOPATH/src/github.com/hashicorp/terraform/vendor/github.com/zclconf/go-cty/cty/types_to_register.go`
2. Comment out the following import statements:
```
 "math/big"
 "github.com/zclconf/go-cty/cty/set"
```
3. Change `func init()` to look like this:
```
 func init() {
   InternalTypesToRegister = []interface{}{
 //    primitiveType{},
 //    typeList{},
 //    typeMap{},
 //    typeObject{},
 //    typeSet{},
 //    setRules{},
 //    set.Set{},
 //    typeTuple{},
 //    big.Float{},
 //    capsuleType{},
     []interface{}(nil),
     map[string]interface{}(nil),
 45   }
 ```
 4. Rebuild laforge using `go build github.com/gen0cide/laforge/cmd/laforge`
