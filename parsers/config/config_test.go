package config

import (
	"testing"
)

type TestSetting struct {
	Test TestStruct
}

type TestStruct struct {
	Value0 int
	Value1 string
}

func Test0(t *testing.T) {
	testconf := "{\"test\":{\"value0\":1,\"value1\":\"something\"}}\r\n#this is a comment\r\n"
	result := TestSetting{}
	err := Parse(testconf, &result)
	t.Log(result)
	t.Log(err)
	if result.Test.Value0 != 1 || result.Test.Value1 != "something" || err != nil {
		t.Error("Error parsing:", err)
	}
}
