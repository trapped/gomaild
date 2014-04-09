package message

import (
	"strings"
	"testing"
)

func Test0Headers(t *testing.T) {
	text0 := "From: Someone\r\nTo: Someone else\r\n\r\nThis is the body.\r\nLine two."
	text1 := "This is the body.\r\nLine two."
	text2 := "\r\nThis is the body.\r\nLine two."

	msg0, msg1, msg2 := Message{Text: text0}, Message{Text: text1}, Message{Text: text2}

	headers0, err0 := Headers(msg0)
	headers1, err1 := Headers(msg1)
	headers2, err2 := Headers(msg2)

	if err0 != nil {
		t.Log("Error getting headers of message #0:\r\n" + err0.Error())
	}
	if err1 != nil {
		t.Log("Error getting headers of message #1:\r\n" + err1.Error())
	}
	if err2 != nil {
		t.Log("Error getting headers of message #2:\r\n" + err2.Error())
	}

	t.Log("Headers of message #0:\r\n" + headers0)
	t.Log("Headers of message #1:\r\n" + headers1)
	t.Log("Headers of message #2:\r\n" + headers2)
}

func Test1AppendHeaders(t *testing.T) {
	text0 := "From: Someone\r\nTo: Someone else\r\n\r\nThis is the body.\r\nLine two."
	text1 := "This is the body.\r\nLine two."
	text2 := "\r\nThis is the body.\r\nLine two."

	msg0, msg1, msg2 := Message{Text: text0}, Message{Text: text1}, Message{Text: text2}

	headers0, err0 := Headers(msg0)
	headers1, err1 := Headers(msg1)
	headers2, err2 := Headers(msg2)

	if err0 != nil {
		t.Log("Error getting headers of message #0:\r\n" + err0.Error())
	}
	if err1 != nil {
		t.Log("Error getting headers of message #1:\r\n" + err1.Error())
	}
	if err2 != nil {
		t.Log("Error getting headers of message #2:\r\n" + err2.Error())
	}

	body0, err0 := Body(msg0)
	body1, err1 := Body(msg1)
	body2, err2 := Body(msg2)

	if err0 != nil {
		t.Log("Error getting body of message #0:\r\n" + err0.Error())
	}
	if err1 != nil {
		t.Log("Error getting body of message #1:\r\n" + err1.Error())
	}
	if err2 != nil {
		t.Log("Error getting body of message #2:\r\n" + err2.Error())
	}

	t.Log("Body of message #0:\r\n" + body0)
	t.Log("Body of message #1:\r\n" + body1)
	t.Log("Body of message #2:\r\n" + body2)

	newheaders0 := append(strings.Split(headers0, "\r\n"), "X-Test: Test 0")
	newheaders1 := append(strings.Split(headers1, "\r\n"), "X-Test: Test 1")
	newheaders2 := append(strings.Split(headers2, "\r\n"), "X-Test: Test 2")

	newfull0 := strings.Join(newheaders0, "\r\n") + "\r\n\r\n" + body0
	newfull1 := strings.Join(newheaders1, "\r\n") + "\r\n\r\n" + body1
	newfull2 := strings.Join(newheaders2, "\r\n") + "\r\n\r\n" + body2

	t.Log("Full message #0:\r\n" + newfull0)
	t.Log("Full message #1:\r\n" + newfull1)
	t.Log("Full message #2:\r\n" + newfull2)
}
