JsonX - JSON extended for configuration files and human use
===========================================================


I really like JSON but it has a few problems.  JsonX is a system that is intended to address these problems.
JsonX is an extended JSON language specifically designed for:

1. Humans to create and edit
2. Configuration Files
3. More readable syntax
4. Syntactic error recovery
5. Validation of input

Like most software this was developed to solve a specific set of pain points.

1. An open source project that allowed syntactically invalid input to result in an irreversible configuration.  Data was sent to browsers that is both incorrect and will be cached for 10 years.
2. Lack of comments inside JSON.   I am just not fond of languages that lack the ability to add clear documentation.
3. Lack of decent error messages.  Syntax errors and the most common use of JSON in Go (the built in package) leads to the error message "You have a syntax error" without any meaningful indication
of where that syntax error may lie.  The last time I used a software language like that was Fortran 66.
4. The failure of JSON parsers to actually follow the JSON specification.  Yes Go's JSON parser has some irritating errors just like most other JSON parsers.
5. The lack of any "time" based data in JSON.

Examples
--------------

1. ./example/simple - A simple example with of reading in a .jsonx file.
2. ./example/ex2 - An example that uses gfDefauilt and isSet with ints
3. ./example/ex3 - The ex2 example modified for strings with input demonstrating different quotes.

Quick Overview
--------------

In Go code

```go

	import (
		"www.2c-why.com/JsonX"
	)

```

Then to read in a JsonX file and unmarsal it:

```go
	fn := "file.jsonx"
	meta, err := JsonX.UnmarshalFile(fn, &TheDataConverted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error in reading or unmarshaling the file %s as JsonX: %s", fn, err )
		os.Exit(1)
	}
	_ = meta // needed to validate and check for missing values!
```

Exampel 1
--------------

Read a configuration file.

```go

type ClsCfgType struct {
	Id    string `gfJsonX:"id"`
	Title string `gfJsonX:"title"`
	Desc  string `gfJsonX:"desc"`
}

func ReadClsCfg(cfg *CfgType, fn string) (clsCfg ClsCfgType, err error) {
	In, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for cls-config, error=%s\n", fn, err)
		return
	}

	meta, err := JsonX.Unmarshal(fn, In, &clsCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "JxCli: Error returned from Unmarshal: %s\n", err)
		return
	}

	err = JsonX.ValidateValues(&clsCfg, meta, "", "", "")
	if err != nil {
		// print out some useful error message at this point includeing "err"
		return
	}

	// Test ValidateRequired - both success and fail, nested, in arrays
	errReq := JsonX.ValidateRequired(&clsCfg, meta)
	if errReq != nil {
		// print out some useful error message including errReq
		fmt.Printf("err=%s meta=%s\n", err, JsonX.SVarI(meta))
		msg, err := JsonX.ErrorSummary("text", &clsCfg, meta)
		if err != nil {
			fmt.Printf("%s\n", msg)
			return clsCfg, err
		}
		return clsCfg, errReq
	}

	return
}

```

The configuration file

```javascript
{
	title: "An introduction to programming using Go"
	desc: "A beginning programming class using Go(golang)."
}
```

Note that "tags" inside the has need not be quoted and that commons are optional.



Exampel 2
--------------

Comments
--------------

Comments start with a `//` to end of line or a `/*` until the next `*/`.
You can configure the multi-line comments to nest. 


File inclusion
-----------------

You can include files from the .jx (JsonX) files.  For Example:

```javascript
{
	IncC: {{ __include__ ./dir/a.file }}
}
```

Where the directory `./dir` contains a.file.  the `{{...}}` will be replaced with the contents of the include file.
Note: The paths are relative to where the file is, so `./` for the top file, then if a.file includes another file,
say b.file it will be relative to the location of a.file. 

File inclusion is processed at a low level, below the level of the scanner.  This means that you should be able
to include files at any location in the input.


Buitin Functions
-----------------

From: scan.go:

```
	rv.Funcs["__include__"] = fxInclude
	rv.Funcs["__sinclude__"] = fxInclude
	rv.Funcs["__require__"] = fxInclude
	rv.Funcs["__srequire__"] = fxInclude
	rv.Funcs["__file_name__"] = fxFileName
	rv.Funcs["__line_no__"] = fxLineNo
	rv.Funcs["__col_pos__"] = fxColPos
	rv.Funcs["__now__"] = fxNow // the current date time samp
	rv.Funcs["__q__"] = fxQ     // return a quoted empty string
	rv.Funcs["__start_end_marker__"] = fxStartEndMarker
	rv.Funcs["__comment__"] = fxComment
	rv.Funcs["__comment_nest__"] = fxCommentNest
	rv.Funcs["__env__"] = fxEnv
```

TODO - document each of the builtin functions

Template processing
--------------------

Pulling Defaults/Values from Environemnt
-----------------------------------------

```
{{ __env__ Name Name Name }}
```

Pulls in values from the environment so that you can substitute them into your
input.

```bash

$ export DB_PASS=mypassword

```

Then

```
{{ __env__ DB_PASS }}
```

Will substitute in the database password, `mypassword`.  If you need to substitute
this into a string then use """ quotes.   See JxCli/testdata/test0021.jx for example.


Default From Environment
--------------------------

You can set the default value for a struct from the environment also.

```go
type SomeStruct struct {
	Field string `gfDefaultEnv:"DB_PASS"`
}
```

Will pull `DB_PASS` from the environment and use the value of that for the default.


Pulling Defaults/Values from Redis 
-------------------------------------

Setting Default Values Using Tags
---------------------------------

A number of tags can be used to set default values inside structures.  These are:
`gfDefault`, `gfDefaultEnv`, `gfDefaultFromKey`, `gfAlloc`.

For example:

```go
	type Example01 struct {
		AnInt	int		`gfDefault:"22"`
		AString	string	`gfDefault:"yep"`
		ABool	bool	`gfDefault:"true"`
		AFloat	float	`gfDefault:"2.2"`
	}
	var example01 Example01
	...
	meta := make(map[string]MetaInfo)
	err := SetDefaults(&example01, meta, "", "", "")
```

Will set `AnInt` to a value of `22`, `AString` to `yep`, `ABool` to `true` and `AFloat` to `2.2`. At first glance this
seems like a difficult and complex way to set default values. However this allows for setting of default values wen the
data is read in from a JsonX file. 

To pull a value from the environment:

```go
	type Example02 struct {
		DatabasePassword	string		`gfDefaultEnv:"DB_PASS"`
	}
	...
	var example02 Example02
	...
	meta := make(map[string]MetaInfo)
	meta, err := JsonX.UnmarshalFile(fn, &example02)
	...
```


Reading in JsonX Data
---------------------

TODO

JsonX Tags
--------------------------

The following tags are supported.

| Tag                   | Description                                                               |
| :-------------------- | :------------------------------------------------------------------------ |
| gfDefault             | Default value for field.                                                  |
| gfDefaultEnv          | Name of environment variable to use for default if not set.               |
| gfDefaultFromKey      | Key passed to PullFromDefault to get default value if not set             |
| gfAlloc               | Allocate a minimum of N items in a slice.                                 |
| gfPromptPassword      | Prompt stdin for a password to fill in field.                             |
| gfPromp               | Prompt stdin for a value to fill in field.                                |
| gfType                | A validation type, int, money, email, url, filepath, filename, fileexists |
| gfMinValue            | Inclusive minimum value for field, int, float or string.                  |
| gfMaxValue            | Inclusive maximum value for field, int, float or string.                  |
| gfListValue           | Value must come from list supplied.                                       |
| gfValieRE             | Regular expression must match value to be valid.                          |
| gfIgnoreDefault       | If a default value is used then validation will not be applied.           |
| gfRequired            | Must be supplied, can not be left empty.                                  |
| gfTag                 | Match name for field.                                                     |
| gfNoSet               | If "-" then will not be set. Like 1st item in `json` tag.                 |








