package message

import (
	"reflect"
	"testing"
)

func TestEdgeCases(t *testing.T) {
	t.Run("empty message", func(t *testing.T) {
		_, err := ParseMessage("")
		if err == nil || err != InvalidInput {
			t.Errorf("Expected error to be %q, got %q", InvalidInput, err)
		}
	})

	t.Run("empty string with carriage return and line feed", func(t *testing.T) {
		_, err := ParseMessage("\r\n")
		if err == nil || err != InvalidInput {
			t.Errorf("Expected error to be %q, got %q", InvalidInput, err)
		}
	})

	t.Run("missing carriage return or line feed at the end", func(t *testing.T) {
		_, err := ParseMessage("input not properly ended")
		if err == nil || err != MissingCRLF {
			t.Errorf("Expected error to be %q, got %q", MissingCRLF, err)
		}
	})

	t.Run("params cannot be more than 14", func(t *testing.T) {
		_, err := ParseMessage("PRIVMSG a b c d e f g h i j k l m n o\r\n")
		if err == nil || err != InvalidParam {
			t.Errorf("Expected error to be %q, got %q", InvalidParam, err)
		}
	})

	t.Run("param cannot be more than 14 characters long", func(t *testing.T) {
		_, err := ParseMessage("PRIVMSG abcdefghijklmno\r\n")
		if err == nil || err != InvalidParam {
			t.Errorf("Expected error to be %q, got %q", InvalidParam, err)
		}
	})
}

func TestParseMessage(t *testing.T) {
	testCases := []struct{
		Name     string
		Input    string
		Expected *Message
	} {
		{
			Name: "Test NO tags, prefix & trailing",
			Input: "PRIVMSG #test hello\r\n",
			Expected: &Message{
				Tags: []string{},
				Prefix: "",
				Command: "PRIVMSG",
				Params: []string{"#test", "hello"},
			},
		},
		{
			Name: "Test NO prefix & trailing",
			Input: "@a=b PRIVMSG #test hello\r\n",
			Expected: &Message{
				Tags: []string{"a=b"},
				Prefix: "",
				Command: "PRIVMSG",
				Params: []string{"#test", "hello"},
			},
		},
		{
			Name: "Test with multiple tags",
			Input: "@a=b;c;d=e;url=http://example.com PRIVMSG #test hello\r\n",
			Expected: &Message{
				Tags: []string{"a=b", "c", "d=e", "url=http://example.com"},
				Prefix: "",
				Command: "PRIVMSG",
				Params: []string{"#test", "hello"},
			},
		},
		{
			Name: "Test prefix",
				Input: "@a=b;c;d=e;url=http://example.com :irc.example.chat PRIVMSG #test hello\r\n",
			Expected: &Message{
				Tags: []string{"a=b", "c", "d=e", "url=http://example.com"},
				Prefix: "irc.example.chat",
				Command: "PRIVMSG",
				Params: []string{"#test", "hello"},
			},
		},
		{
			Name: "Test numeric command",
				Input: "@a=b;c;d=e;url=http://example.com :irc.example.chat 254 #test hello\r\n",
			Expected: &Message{
				Tags: []string{"a=b", "c", "d=e", "url=http://example.com"},
				Prefix: "irc.example.chat",
				Command: "254",
				Params: []string{"#test", "hello"},
			},
		},
		{
			Name: "Test trailing",
				Input: "@a=b;c;d=e;url=http://example.com :irc.example.chat 254 #test hello :this is the trailing part of the message\r\n",
			Expected: &Message{
				Tags: []string{"a=b", "c", "d=e", "url=http://example.com"},
				Prefix: "irc.example.chat",
				Command: "254",
				Params: []string{"#test", "hello", "this is the trailing part of the message"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func (t *testing.T) {
			msg, err := ParseMessage(testCase.Input)
			if err != nil {
				t.Fatalf("Unexepected error: %q", err)
			}
			if !reflect.DeepEqual(msg, testCase.Expected) {
				t.Fatalf("Expected %q, got %q", testCase.Expected, msg)
			}
		})
	}
}
