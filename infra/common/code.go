package common

const (
	ERR_CODE_OK        int32 = 0
	ERR_CODE_PARAM_ERR int32 = 1
	ERR_CODE_SYS_ERR   int32 = 2
	ERR_CODE_TOKEN_ERR int32 = 3
)

var ErrMsg = map[int32]string{
	ERR_CODE_OK:        "ok",
	ERR_CODE_PARAM_ERR: "param err",
	ERR_CODE_SYS_ERR:   "sys err",
	ERR_CODE_TOKEN_ERR: "token err",
}
