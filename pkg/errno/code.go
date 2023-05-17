package errno

// 定义错误类型
// 错误类型通常包含两部分: Code 部分，用来唯一标识一个错误；Message 部分，用来展示错误信息，这部分错误信息通常供前端直接展示
// 错误代码说明
// 错误代码包含3部分：
// 1. 服务级别代码：1 为系统级错误；2 为普通错误(通常是由用户非法操作引起的)
// 2. 服务模块代码：
// 3. 具体错误代码：
// code = 0 说明是正确返回，code > 0 说明是错误返回

var (
	// 一般错误
	OK                  = &Errno{Code: 0, Message: "OK"}                        // 运行正常
	InternalServerError = &Errno{Code: 10000, Message: "Internal server error"} // 内部服务错误
	ErrorBind           = &Errno{Code: 10001, Message: "Request binding error"} // 参数绑定错误

	ErrorValidation       = &Errno{Code: 20000, Message: "Validation failed"} // 参数非法
	ErrorDatabase         = &Errno{Code: 20001, Message: "Database error"}    // 数据库错误
	ErrorToken            = &Errno{Code: 20002, Message: "Signing the JSON web token error"}
	ErrorMissingParameter = &Errno{Code: 20003, Message: "Missing parameter"}       //缺少参数
	ErrorParameterParsing = &Errno{Code: 20004, Message: "Parameter parsing error"} // 参数解析错误

	// user errors
	ErrorEncrypt        = &Errno{Code: 20100, Message: "Error occurred while encrypting the user password"} // 用户密码加盐出错
	ErrorLogin          = &Errno{Code: 20101, Message: "Wrong username or password"}                        // 登录过程用户名密码错误
	ErrorTokenInvalid   = &Errno{Code: 20102, Message: "Unauthorized"}                                      // 未授权
	ErrorUsernameRepeat = &Errno{Code: 20103, Message: "Username Repeat"}                                   // 注册用户名重复
	// ErrorUserNotFound      = &Errno{Code: 20101, Message: "User not found"}
	// ErrorPasswordIncorrect = &Errno{Code: 20102, Message: "Password incorrect"}

	ErrorAlgorithmCount = &Errno{Code: 20200, Message: "Error Algorithm Count"} // 定位算法执行错误

	ErrorRecordNotFound = &Errno{Code: 20300, Message: "Record not found"} // 数据库查询记录不存在
)
