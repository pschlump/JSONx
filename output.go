package JsonX

// Copyright (C) Philip Schlump, 2014-2017

//
// Output
//
// Uses the line with configuration so that in JSON we get:
//		{
//			"abc":"def",
//			"ghi":"jkl"
//		}
// In JsonX we get for a bottom level item
//		{ "abc":"def", "ghi":"jkl" }
// Until the line is full, then a wrap if OutputInJSON is true and
//		{ abc:"def", ghi:"jkl" }
// if the flag is false.  IDs will only be quoted if they contain non-id characters (like blank for example).
//
