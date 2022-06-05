package accept

import "testing"

func TestMatches(t *testing.T) {
	tests := []struct {
		Matcher  string
		Absolute string
		Result   bool
	}{
		{
			Matcher:  "*/*",
			Absolute: "any/thing",
			Result:   true,
		},
		{"text/*", "text/html", true},
		{"text/html", "text/html", true},
		{"text/*", "application/json", false},
		{"text/plain", "text/html", false},
		{"application/json", "application/ld+json", true},
		{"application/activity+json", "application/json", false},
	}
	for i, test := range tests {
		if Matches(test.Matcher, test.Absolute) != test.Result {
			t.Errorf("Test %d failed. %v %v %v", i, test.Matcher, test.Absolute, test.Result)
		}
	}
}

func TestParseheader(t *testing.T) {
	test1 := "text/*, text/plain, text/plain;format=flowed, */*"
	parsed := ParseHeader(test1)
	if parsed[0].MediaRange != "text/*" {
		t.Error("Test1 failed", test1, parsed)
	}
	if parsed[2].MediaRange != "text/plain" {
		t.Error("Test1 failed", test1, parsed)
	}

	test2 := "text/*;q=0.3, text/html;q=0.7, text/html;level=1,text/html;level=2;q=0.4, */*;q=0.5"
	parsed2 := ParseHeader(test2)
	if parsed2[0].MediaRange != "text/html" {
		t.Error("Test2 failed", test2, "\n", parsed2)
	}
	if parsed2[2].MediaRange != "*/*" {
		t.Error("Test2 failed", test2, "\n", parsed2)
	}
}

func TestServeType(t *testing.T) {
	test1 := "audio/*; q=0.2, audio/basic"
	servable := []string{"audio/mp3", "audio/wav", "audio/basic"}
	if r := ServeType(servable, test1); r != "audio/basic" {
		t.Error("Test failed", r)
	}
}
