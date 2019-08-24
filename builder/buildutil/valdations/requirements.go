package validations

import (
	"os"
	"os/exec"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/iancoleman/strcase"

	"github.com/gen0cide/laforge/core"
	"github.com/gen0cide/laforge/core/cli"
)

// Validations is suppose to be a sequence of validations that must pass for a builder
type Validations []Requirement

// Requirement defines a named requirement that must be met for a builder to continue
type Requirement struct {
	Name       string
	Resolution string
	Check      Check
}

// Check is a type alias to describe a method that validates a laforge context
type Check func(base *core.Laforge) bool

// Not is the logical negation of a Check
func Not(f Check) Check {
	return func(base *core.Laforge) bool {
		return !f(base)
	}
}

// All returns a meta-Check that requires all supplied Checks evaluate to true
func All(funcs ...Check) Check {
	return AtLeastN(len(funcs), funcs...)
}

// And creates a meta Check for requiring a logical AND between supplied Checks
func And(a, b Check) Check {
	return All(a, b)
}

// Any creates a meta Check for any possible hits against a set of Checks
func Any(funcs ...Check) Check {
	return AtLeastN(1, funcs...)
}

// Or creates a meta Check for performing a logical OR on two Checks
func Or(a, b Check) Check {
	return Any(a, b)
}

// AtLeastN creates an at least n rule against a set of Checks
func AtLeastN(n int, funcs ...Check) Check {
	if n < 1 {
		n = 1
	}
	if n > len(funcs) {
		n = len(funcs)
	}
	return func(base *core.Laforge) bool {
		passes, fails := 0, 0
		for _, f := range funcs {
			if f(base) {
				passes++
			} else {
				fails++
			}
			if len(funcs)-fails < n {
				return false
			}
			if passes >= n {
				return true
			}
		}
		return false
	}
}

// ExistsInPath checks to see if a command line tool is installed and in the current user's path
func ExistsInPath(progname string) Check {
	return func(base *core.Laforge) bool {
		_, err := exec.LookPath(progname)
		return err == nil
	}
}

// FieldNotEmpty checks a type in the context state for a nil/zero value field and fails if it is so.
func FieldNotEmpty(obj interface{}, fieldname string) Check {
	camName := strcase.ToCamel(fieldname)
	switch v := obj.(type) {
	case core.Competition:
		return func(base *core.Laforge) bool {
			if validation.IsEmpty(base) {
				cli.Logger.Errorf("base state was empty")
				return false
			}
			if validation.IsEmpty(base.CurrentCompetition) {
				cli.Logger.Errorf("competition state was empty")
				return false
			}
			compVal := reflect.Indirect(reflect.ValueOf(base.CurrentCompetition))
			if compVal.Kind() != reflect.Struct {
				cli.Logger.Errorf("Competition has failed a validation error: base.Competition was not of type struct")
				return false
			}
			fieldTest := compVal.FieldByName(camName)
			if !fieldTest.IsValid() {
				cli.Logger.Errorf("%s does not have a field named %s!", "competition", camName)
				os.Exit(1)
				return false
			}
			testField := fieldTest.Interface()
			if validation.IsEmpty(testField) {
				cli.Logger.Errorf("Competition has field a validation error: field %s was empty", camName)
				return false
			}
			return true
		}
	case core.DNS:
		return func(base *core.Laforge) bool {
			if validation.IsEmpty(base) {
				cli.Logger.Errorf("base state was empty")
				return false
			}
			if validation.IsEmpty(base.CurrentCompetition) {
				cli.Logger.Errorf("competition state was empty")
				return false
			}
			if validation.IsEmpty(base.CurrentCompetition.DNS) {
				cli.Logger.Errorf("dns state was empty")
				return false
			}
			compVal := reflect.Indirect(reflect.ValueOf(base.CurrentCompetition.DNS))
			if compVal.Kind() != reflect.Struct {
				cli.Logger.Errorf("DNS has failed a validation error: base.CurrentCompetition.DNS was not of type struct")
				return false
			}
			fieldTest := compVal.FieldByName(camName)
			if !fieldTest.IsValid() {
				cli.Logger.Errorf("%s does not have a field named %s!", "dns", camName)
				os.Exit(1)
				return false
			}
			testField := fieldTest.Interface()
			if validation.IsEmpty(testField) {
				cli.Logger.Errorf("DNS has field a validation error: field %s was empty", camName)
				return false
			}
			return true
		}
	case core.Environment:
		return func(base *core.Laforge) bool {
			if validation.IsEmpty(base) {
				cli.Logger.Errorf("base state was empty")
				return false
			}
			if validation.IsEmpty(base.CurrentEnv) {
				cli.Logger.Errorf("Environment state was empty")
				return false
			}
			compVal := reflect.Indirect(reflect.ValueOf(base.CurrentEnv))
			if compVal.Kind() != reflect.Struct {
				cli.Logger.Errorf("Environment has failed a validation error: base.Environment was not of type struct")
				return false
			}
			fieldTest := compVal.FieldByName(camName)
			if !fieldTest.IsValid() {
				cli.Logger.Errorf("%s does not have a field named %s!", "environment", camName)
				os.Exit(1)
				return false
			}
			testField := fieldTest.Interface()
			if validation.IsEmpty(testField) {
				cli.Logger.Errorf("Environment has field a validation error: field %s was empty", camName)
				return false
			}
			return true
		}
	case core.Host:
		return func(base *core.Laforge) bool {
			for n, o := range base.Hosts {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "host", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "host", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "host", n, camName)
					return false
				}
			}
			return true
		}
	case core.Command:
		return func(base *core.Laforge) bool {
			for n, o := range base.Commands {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "command", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "command", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "command", n, camName)
					return false
				}
			}
			return true
		}
	case core.Identity:
		return func(base *core.Laforge) bool {
			for n, o := range base.Identities {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "identity", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "identity", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "identity", n, camName)
					return false
				}
			}
			return true
		}
	case core.Network:
		return func(base *core.Laforge) bool {
			for n, o := range base.Networks {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "network", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "network", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "network", n, camName)
					return false
				}
			}
			return true
		}
	case core.RemoteFile:
		return func(base *core.Laforge) bool {
			for n, o := range base.RemoteFiles {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "remote_file", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "remote_file", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "remote_file", n, camName)
					return false
				}
			}
			return true
		}
	case core.Script:
		return func(base *core.Laforge) bool {
			for n, o := range base.Scripts {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "script", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "script", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if validation.IsEmpty(testField) {
					cli.Logger.Errorf("%s %s has field a validation error: field %s was empty", "script", n, camName)
					return false
				}
			}
			return true
		}
	default:
		return func(base *core.Laforge) bool {
			cli.Logger.Errorf("Invalid type %T passed for validation", v)
			return false
		}
	}
}

// FieldEquals allows for comparison of various fields within the Laforge state during builder validation
func FieldEquals(obj interface{}, fieldname string, equals interface{}) Check {
	camName := strcase.ToCamel(fieldname)
	switch v := obj.(type) {
	case core.Competition:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentCompetition != nil {
				compVal := reflect.Indirect(reflect.ValueOf(base.CurrentCompetition))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("Competition has failed a validation: base.Competition was not of type struct")
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "competition", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("competition field %s (%T) was not equal to %v (%T)", camName, testField, equals, equals)
				return false
			}
			cli.Logger.Errorf("Competition has failed a validation: base or base.Competition was nil")
			return false
		}
	case core.DNS:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentCompetition != nil && base.CurrentCompetition.DNS != nil {
				compVal := reflect.Indirect(reflect.ValueOf(base.CurrentCompetition.DNS))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("Competition has failed a validation: base.Competition was not of type struct")
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "dns", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("DNS field %s (%T) was not equal to %v (%T)", camName, testField, equals, equals)
				return false
			}
			cli.Logger.Errorf("Competition has failed a validation: base or base.Competition or base.Competition.DNS was nil")
			return false
		}
	case core.Environment:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentEnv != nil {
				compVal := reflect.Indirect(reflect.ValueOf(base.CurrentEnv))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("Environment has failed a validation: base.Environment was not of type struct")
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "environment", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("environment field %s (%T) was not equal to %v (%T)", camName, testField, equals, equals)
				return false
			}
			cli.Logger.Errorf("Environment has failed a validation: base or base.Environment was nil")
			return false
		}
	case core.Host:
		return func(base *core.Laforge) bool {
			for n, o := range base.Hosts {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "host", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "host", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "host", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	case core.Command:
		return func(base *core.Laforge) bool {
			for n, o := range base.Commands {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "command", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "command", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "command", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	case core.Identity:
		return func(base *core.Laforge) bool {
			for n, o := range base.Identities {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "identity", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "identity", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "identity", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	case core.Network:
		return func(base *core.Laforge) bool {
			for n, o := range base.Networks {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "networks", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "network", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "networks", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	case core.RemoteFile:
		return func(base *core.Laforge) bool {
			for n, o := range base.RemoteFiles {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "remote_file", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "remote_file", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "remote_file", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	case core.Script:
		return func(base *core.Laforge) bool {
			for n, o := range base.Scripts {
				compVal := reflect.Indirect(reflect.ValueOf(o))
				if compVal.Kind() != reflect.Struct {
					cli.Logger.Errorf("%s %s has failed a validation: was not of type struct (found %s)", "script", n, compVal.Kind().String())
					return false
				}
				fieldTest := compVal.FieldByName(camName)
				if !fieldTest.IsValid() {
					cli.Logger.Errorf("%s does not have a field named %s!", "scripts", camName)
					os.Exit(1)
					return false
				}
				testField := fieldTest.Interface()
				if reflect.DeepEqual(testField, equals) {
					return true
				}
				cli.Logger.Errorf("%s %s has field mismatch: field %s (type=%T,val=%v) was not equal to %v (%T)", "script", n, camName, testField, testField, equals, equals)
				return false
			}
			return true
		}
	default:
		return func(base *core.Laforge) bool {
			cli.Logger.Errorf("Invalid type %T passed for validation", v)
			return false
		}
	}
}

// MapHasKey is a small helper function to check if a map has an key (key)
func MapHasKey(key string, m map[string]string) bool {
	if _, found := m[key]; found {
		return true
	}
	return false
}

// HasConfigKey checks configurable types (Competition, DNS, Environment) for configuration values.
func HasConfigKey(obj interface{}, key string) Check {
	switch v := obj.(type) {
	case core.Competition:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentCompetition != nil && MapHasKey(key, base.CurrentCompetition.Config) {
				return true
			}
			cli.Logger.Errorf("Competition has failed a validation: config parameter %s was not defined", key)
			return false
		}
	case core.DNS:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentCompetition != nil && base.CurrentCompetition.DNS != nil && MapHasKey(key, base.CurrentCompetition.DNS.Config) {
				return true
			}
			cli.Logger.Errorf("DNS has failed a validation: config parameter %s was not defined", key)
			return false
		}
	case core.Environment:
		return func(base *core.Laforge) bool {
			if base != nil && base.CurrentEnv != nil && MapHasKey(key, base.CurrentEnv.Config) {
				return true
			}
			cli.Logger.Errorf("Environment has failed a validation: config parameter %s was not defined", key)
			return false
		}
	default:
		return func(base *core.Laforge) bool {
			cli.Logger.Errorf("Invalid type %T passed for validation", v)
			return false
		}
	}
}

// HasVarDefined checks types supporting variable assignment to see if they have a variable of a specific name assigned.
// valid object classes: Command, Host, Identity, Network, RemoteFile, Script
func HasVarDefined(obj interface{}, varname string) Check {
	switch obj.(type) {
	case core.Host:
		return func(base *core.Laforge) bool {
			for n, o := range base.CurrentEnv.IncludedHosts {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Host %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	case core.Identity:
		return func(base *core.Laforge) bool {
			for n, o := range base.Identities {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Identity object %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	case core.Network:
		return func(base *core.Laforge) bool {
			for n, o := range base.Networks {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Network object %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	case core.RemoteFile:
		return func(base *core.Laforge) bool {
			for n, o := range base.RemoteFiles {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Remote File object %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	case core.Script:
		return func(base *core.Laforge) bool {
			for n, o := range base.Scripts {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Script object %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	case core.Command:
		return func(base *core.Laforge) bool {
			for n, o := range base.Commands {
				if MapHasKey(varname, o.Vars) {
					continue
				}
				cli.Logger.Errorf("Command object %s has failed a validation: var %s was not defined", n, varname)
				return false
			}
			return true
		}
	default:
		return func(base *core.Laforge) bool {
			return false
		}
	}
}
