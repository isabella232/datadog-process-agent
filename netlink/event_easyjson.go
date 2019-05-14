// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package netlink

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF642ad3eDecodeGithubComDataDogDatadogProcessAgentNetlink(in *jlexer.Lexer, out *IPTranslation) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "repl_src_ip":
			out.ReplSrcIP = string(in.String())
		case "repl_dst_ip":
			out.ReplDstIP = string(in.String())
		case "repl_src_port":
			out.ReplSrcPort = uint16(in.Uint16())
		case "repl_dst_port":
			out.ReplDstPort = uint16(in.Uint16())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF642ad3eEncodeGithubComDataDogDatadogProcessAgentNetlink(out *jwriter.Writer, in IPTranslation) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"repl_src_ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ReplSrcIP))
	}
	{
		const prefix string = ",\"repl_dst_ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ReplDstIP))
	}
	{
		const prefix string = ",\"repl_src_port\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint16(uint16(in.ReplSrcPort))
	}
	{
		const prefix string = ",\"repl_dst_port\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint16(uint16(in.ReplDstPort))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v IPTranslation) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeGithubComDataDogDatadogProcessAgentNetlink(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IPTranslation) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeGithubComDataDogDatadogProcessAgentNetlink(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *IPTranslation) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeGithubComDataDogDatadogProcessAgentNetlink(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IPTranslation) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeGithubComDataDogDatadogProcessAgentNetlink(l, v)
}
