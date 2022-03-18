package stripjsonc

import (
	"os"
	"testing"
	"unicode/utf8"
)

func TestStripJSONC(t *testing.T) {
	var input = `// this is a single line comment
{ "foo //": 123, // another single line comment
  /* 
  * block comments
  * 👍
  *
/ * /
  can contain whatever you "want"
  */ "b😃\"\\ar/*": /* inline block comments are ok */"*/baz"
}
`
	var expected = `                                
{ "foo //": 123,                               
     
                  
        
   
     
                                 
     "b😃\"\\ar/*":                                   "*/baz"
}
`
	output := StripJSONCString(input)
	if output != expected {
		t.Errorf("Output %s jsonc does not match expected %s", output, expected)
	}
}

func FuzzStripJSONC(f *testing.F) {
	testcases := []string{
		"// just a comment",
		"[\"f//oo\", /*\n block */ 2]",
		"/* unterminated block",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		res := StripJSONCString(orig)
		if utf8.ValidString(orig) && !utf8.ValidString(res) {
			t.Errorf("StripJSONC produced invalid UTF-8 string %q", res)
		}
	})
}

func BenchmarkStripJSONCStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("data/large.jsonc")
		StripJSONCStream(f, os.Stdout)
		f.Close()
	}
}
