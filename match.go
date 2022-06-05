package accept

import (
	"fmt"
	"mime"
	"sort"
	"strconv"
	"strings"
)

type AcceptExt struct {
	MediaRange string
	Params     map[string]string
}

func ServeType(servable []string, acceptheader string) string {
	requested := ParseHeader(acceptheader)
	for _, acpt := range requested {
		for _, s := range servable {
			if Matches(acpt.MediaRange, s) {
				return s
			}
		}
	}
	return ""
}

func ParseHeader(header string) (accepts []AcceptExt) {
	parts := strings.Split(header, ",")
	for _, t := range parts {
		mediatype, params, err := mime.ParseMediaType(t)
		if err == nil {
			accepts = append(accepts, AcceptExt{mediatype, params})
		} else {
			fmt.Println("Error", err)
		}
	}
	Sort(accepts)
	return
}

func Matches(matcher string, absolute string) bool {
	if matcher[0] == byte('*') {
		return true
	}
	matcherspl := strings.Split(matcher, "/")
	absolutspl := strings.Split(absolute, "/")
	if matcherspl[0] != absolutspl[0] {
		return false
	}
	if len(matcherspl) < 2 || len(absolutspl) < 2 {
		return false
	}
	if matcherspl[1] == "*" {
		return true
	}
	if matcherspl[1] == absolutspl[1] {
		return true
	}

	// "json" should also match "ld+json" or "activity+json"
	absplit := strings.Split(absolutspl[1], "+")
	if len(absplit) == 1 {
		return false
	}

	if matcherspl[1] == absplit[1] {
		return true
	}
	return false
}

func Sort(items []AcceptExt) {
	sort.Slice(items, func(i, j int) bool {
		iqs, iok := items[i].Params["q"]
		iq, err := strconv.ParseFloat(iqs, 32)
		if !iok || err != nil {
			iq = 1
		}
		jqs, jok := items[j].Params["q"]
		jq, err := strconv.ParseFloat(jqs, 32)
		if !jok || err != nil {
			jq = 1
		}
		return iq > jq
	})
}
