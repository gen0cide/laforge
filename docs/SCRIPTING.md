# Template Scripting Tutorial

There are several places within `laforge` that you have the opportunity to use templates to customize options.

These are implemented using Golang templates. I find their documentation a bit overwhelming, but it's worth giving it a glance over. If you read it and think, "WTF?" - this tutorial is definitely for you.

## Where can I use templates?
Templates can be used in the following locations:

* Setup Scripts (`LF_HOME/scripts/*`)
* Host YAML under `dns_entries` option `name` and `value`

## Templating 101
To use templates, simply use double curly brackets `{{  }}`. Using a shell script as an example:

```
echo "foo"
echo "{{ "bar" }}"
```

This will generate the following shell script:

```
echo "foo"
echo "bar"
```

You do not need to quote the "handle bars" or escape anything inside of them.

Unlike other templating Languages (ERB, Jinja2, etc.), Golang Templates are extremely restrictive around execution. There is only a limited set of functions that can be used. I will outline them below.

## Variable Definition / Function Usage

To define a variable, the following syntax can be used:

```
{{ $foo := printf "%s%d" "bar" 123 }}
echo "{{ $foo }}"
```

This will produce the following shell script:
```

echo "bar123" 
```

The important thing to remember about variable definitions:

* Variables always begin with `$` characters.
* Variables are defined using Golang's `:=` syntax.
* Functions never have parenthesis around arguments list
* Functions do not comma delimit arguments.

### Note about Limiting White Space
When defining variables, you will notice (as in the example above) that the variable declaration created whitespace. If you wish to not have whitespace inserted, use the following syntax:

```
{{ $foo := printf "%s%d" "bar" 123 -}}
echo "{{ $foo }}"
```

Note the ending `-` before the closing `}}`. This tells the template engine to trim whitespace on the end of this definition, thus moving echo up a line.

### printf (your best friend)

You will use `printf` often in templating. If you know the C version, you'll be fine. Look at the Golang docs for a more in depth example and syntax definitions.

## Conditionals

### If Statements

An example:

```
{{ if eq $host.OS "w2k16" }}
  echo "We are on Windows!"
{{ else if eq $host.OS "ubuntu" }}
  echo "We are on Ubuntu!" 
{{ else }}
  echo "We have no idea!"
{{ end }}
``` 

Note the strange syntax of `if FUNC VAL1 VAL2...` instead of what we'd commonly think of `if VAL1 FUNC VAL2`. That will play out often in this tutorial!

Functions that can be used in an if statement:

* `eq` - Returns the boolean truth of `arg1 == arg2`
* `eq` - Returns the boolean truth of `arg1 == arg2`
* `ne` - Returns the boolean truth of `arg1 != arg2`
* `lt` - Returns the boolean truth of `arg1 < arg2`
* `le` - Returns the boolean truth of `arg1 <= arg2`
* `gt` - Returns the boolean truth of `arg1 > arg2`
* `ge` - Returns the boolean truth of `arg1 >= arg2`

### Iterate Maps and Slices
In Golang, Slice is functionally equivalent to an Array while Maps are functionally equivalent to a Hash. There are several situations when you might with to iterate over either data type.

The function you will use for iteration is `range`.

#### Slice Example

```
{{ range $idx, $portNumber := $host.UDPPorts }}
  echo "udp_ports[{{ $idx }}] = {{ $portNumber }}
{{ end }}
```

Here, you'll notice that you have to define two variables in range. The first variable declared is the index offset of the element in the slice, while the second variable is the item in the slice. If you do not need the index offset, you can simply re-write this to:

```
{{ range $_, $portNumber := $host.UDPPorts }}
```

This still will loop through the slice and allow you access to the `$portNumber` object.

#### Map Example

```
{{ range $key, $val := $foo.AnExampleMap }}
  echo "{{ $key }} => {{ $val }}" 
{{ end }}
```

## Exposed Variables

In many of the scripting environments in LaForge, there are a number of variables accessible. Here is a complete list that are available to templates:

* `$.Competition` Competition Object
* `$.Environment` Environment Object
* `$.PodID` Team ID as Integer
* `$.Network` Network Object
* `$.Host` Host Object

Note that these are dynamically defined - you do not have to declare these. They will always reference the scope of the template render. If two hosts on two different networks are using the same `scripts/foo.sh` template, they will have any values rendered according to the scope of the caller.

This allows you to follow DRY principles and write for a broad support of situations.

Please make sure to read the Godocs on the corresponding Object Types and Definitions. They can be found in the `competition` package of this repository.

## Exposed Functions

Besides the build in functions we've already outlined, you cannot typically call Golang functions. I have configured the template engine to embed the following functions for your use.

### N

`N` is a meta function for "counting". In Ruby, this is tribal with the following syntax:

```
10.times do |x|
  puts x
end
```

While this is a normal occurrence, there is no functionally equivalent function in Golang templates. This is where `N` comes in.

You can use `N` to do the same thing as above by doing:

```
{{ range $x, $_ := N 10 }}
  {{ x }}
{{ end }}
```

You give `N` an Integer in a `range` statement, and it will count this for you.
 

### CustomIP

Short helper function to generate valid IP addresses.

```
{{ $ip := CustomIP "192.168.0.0/24" 5 15 }}
```

In the above example, `$ip` will now be a string with the value of `192.168.0.20`

The function is defined as:

```
func CustomIP(cidr string, offset, id int) string {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	newIP := ip.To4()
	lastOctet := offset + id
	newIP[3] = byte(lastOctet)
	return newIP.String()
}
```

### MyIP

```
{{ $myIP := MyIP }}
```

This will trigger a query to an external web service to get your public egress IP and return a string IPv4 representation of it.

### GetUsersForHost

You can use `range` to iterate over the users for a given host:

```
{{ range $_, $user := GetUsersForHost $.Competition $.Host }}
  {{ $user.Email }}
{{ end }}
```

GetUsersForHost takes a `Competition` and `Host` object as arguments and returns a `[]User`.

### GetUsersByOU

Similar to `GetusersForHost`, this function will return a `[]User` for a given OU (represented by a string).

```
{{ range $_, $user := GetUsersByOU "security" }}
  {{ $user.Email }}
{{ end }}
```

## Gotchas

Make sure you familiarize yourself with the `Competition`, `Environment`, `Network`, `Host`, and `User` object definitions. You will use their attributes often.
