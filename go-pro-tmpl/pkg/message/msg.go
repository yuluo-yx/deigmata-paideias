package message

var MsgFlags = map[int]string{
	Success:       "ok",
	Error:         "fail",
	InvalidParams: "args error",
}

// GetMsg Get code mapping message.
func GetMsg(code int) string {

	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}

	return msg
}
