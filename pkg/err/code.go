package errno

// 错误类型通常包含两部分: Code 部分，用来唯一标识一个错误；Message 部分，用来展示错误信息，这部分错误信息通常供前端直接展示
var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error."}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	// user errors
	ErrUserNotFound = &Errno{Code: 20102, Message: "The user was not found."}
)

/*
	错误代码说明
	错误代码包含3部分：
	1. 服务级别代码：1 为系统级错误；2 为普通错误(通常是由用户非法操作引起的)
	2. 服务模块代码：
	3. 具体错误代码：

	code = 0 说明是正确返回，code > 0 说明是错误返回
*/