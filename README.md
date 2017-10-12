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
	meta, err := JsonX.UnmarshalFile(fn, &TheDataConverted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error in reading or unmarshaling the file %s as JsonX: %s", fn, err )
		os.Exit(1)
	}
	_ = meta // need to validate and check for missing values!
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
		// xyzzy
		return
	}

	// Test ValidateRequired - both success and fail, nested, in arrays
	errReq := JsonX.ValidateRequired(&clsCfg, meta)
	if errReq != nil {
		// xyzzy
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

Note that "tags" inside the has need not be quoted and that commans are optional.



Exampel 2
--------------

Comments
--------------

Comments start with a `//` to end of line or a `/*` until the next `*/`.
You can configure the multi-line comments to nest. 

TODO- How to configure.

TODO- How to configure using a built in function - incline.


File inclusion
-----------------

You can include files from the .jx (JsonX) files.  For Example:

```javascript
{
	IncC: {{ __include__ ./dir/a.file }}
}
```

Where the directory ./dir contains a.file.  the {{...}} will be replaced with the contents of the include file.
Note: The paths are relative to where the file is, so ./ for the top file, then if a.file includes another file,
say b.file it will be relative to the location of a.file. 

File inclusion is processed at a low level, below the level of the scanner.  This means that you should be albe
to include files at any location in the input.

TODO - include search path!

TODO - include relative to home directory!

TODO - include relative to a named home directory, `~bob/` for example.

TODO -- Include of "dirctory/Pattern" or "direcoty" - all files in a directory - converted to -- verify this works -- { filename: data, filename2: data2 }

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

	// xyzzy - how do these 2 get set? by user/config
	rv.Data["path"] = []string{"./"}
	rv.Data["envVars"] = []string{}

	rv.Data["fns"] = []string{""}
```

Template processing
--------------------

Pulling Defaults/Values from Environemnt
-----------------------------------------

```
{{ __envVars__ Name Name Name }}
```

Pulls in values from the environment so that you can substitute them into your
input.

```bash

$ export DB_PASS=mypassword

```

Then

```
{{ __envVars__ DB_PASS }}{{.DB_PASS}}
```

Will substitute in the database password, `mypassword`.

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

Data Validation Using Tags
--------------------------

Validation of Interface Data
----------------------------


TODO:
	1. Document example with "is-set" stuff
		set1_test.go:		Mmm    string         `gfJsonX:"mmm,isSet,setField:MmmSet"`
		                    MmmSet bool           `gfJsonX:"-"`
	2. Document use of interface{} - with a JsonX spec for validation and defaults
	3. Document/Implement "output"


Note:
	Options for gfJsonX - "template" indicates substitution for {{.LineNo}} and {{.FileName}} in string value

	Need to have {{.LineNo}} {{.ColPos}} and {{.FileName}} acted on in "scanner"

```
// TODO:
//	1. Add in an external "JsonX" "spec" that specifies the data layout/validation.
//

//
// Tags Are:
//    	gfDefault				The default value
//		gfDefaultEnv			If the speicifed environemnt variable is set then use that instead of gfDefault (gfDefault must be set)
//		gfDefaultFromKey		If this is set and a PullFromDefault function is supplied then call that function to get the value.  Overfides dfDefaultEnv
//		gfIgnore				Ignore setting of deaults for the underlying structure/slice/array - only partially implemnented.
//		gfAlloc					For slice alocate this number of entries and initialize them.  For map allocate a map so you do not have a nil pointer.
//								A Slice 2,10 would be alloacate and initialize 2 with a cap of 10.  For pointers "y" indicates a NIL is to be replace with
//								an pointer to an element.
// Added:
//		gfPromptPassword		Take value from stdin - read as a password without echo.
//		gfPrompt				Take value from stdin
// Note: need to be able to prompt for and validate stuff that is in "Extra" or map[string]interface! -- Meta spec --
//
//
//
//		gfType					(ValidateValues) The type-checking data type, int, money, email, url, filepath, filename, fileexists etc
//		gfMinValue				(ValidateValues) Minimum value
//		gfMaxValue				(ValidateValues) Minimum value
//		gfListValue				(ValidateValues) Value (string) must be one of list of possible values
//		gfValidRE				(ValidateValues) Value will be checked agains a regular expression, match is valid
//		gfIgnoreDefault			(ValidateValues) If value is a default (if mm.SetBy == NotSet) then it will not be validated by above rules.
//
//		gfRequired				(ValidateRequired) Will check that the user supplied a value - not SetBy == NotSet at end
//
//
//		gfTag					(hjson) Tag name for hjson
//		gfNoSet					(hjson) "-" to not set a value, not searialize it. (hjson)
//
// PreCreate, PostCreate - functions to be added -- as an interface/lookup based on type.
//
// Notes:
// 		See: http://intogooglego.blogspot.com/2015/06/day-17-using-reflection-to-write-into.html
// 		See: https://gist.github.com/hvoecking/10772475 -- Looks like a reflection based deep copy
// 		From: http://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces-in-go
// 		http://stackoverflow.com/questions/18091562/how-to-get-underlying-value-from-a-reflect-value-in-golang
```











TODO

isSet v.s. isSetNoDefault - Default name is {Name}IsSet

TODO
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

1. How is the search "PATH" set for include in scan.
2. Listing of included files (tree form listing?)
3. 


http://sosedoff.com/2016/07/16/golang-struct-tags.html -- Validation with Struct Tags

