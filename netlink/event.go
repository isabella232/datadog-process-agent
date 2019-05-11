package netlink

//easyjson:json
type IPTranslation struct {
	ReplSrcIP   string `json:"repl_src_ip"`
	ReplDstIP   string `json:"repl_dst_ip"`
	ReplSrcPort uint16 `json:"repl_src_port"`
	ReplDstPort uint16 `json:"repl_dst_port"`
}
