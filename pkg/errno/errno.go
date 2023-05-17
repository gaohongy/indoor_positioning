package errno

import "fmt"

/* ------------------------------------------------------------------------------ */

// 定义基础错误类型，用于前端显示
type Errno struct {
	Code    int    // 业务响应状态码
	Message string // 业务响应消息
}

// 定义详细错误类型，用于日志信息
type Err struct {
	Code    int    // 业务响应状态码
	Message string // 业务响应消息
	Err     error  // 接口类型，自定义错误
}

/* ------------------------------------------------------------------------------ */
// 实现error这个接口的类型都可以作为一个错误使用，Error 这个方法提供了对错误的描述
// 因此Error和Err下面实现Error()后，都属于一个错误类型，且是自定义类型
// 无论是在控制台还是在日志中输出错误，输出程序都需要知道一个错误应当如何输出，而输出规则就是Error()给出的

// 实现了error接口中声明的Error()
func (err Errno) Error() string {
	return err.Message
}

func (err *Err) Error() string {
	return fmt.Sprintf("error code: %d, message: %s, detail: %s", err.Code, err.Message, err.Err)
}

/* ------------------------------------------------------------------------------ */

// @title	New
// @description	新建自定义错误类型
// @auth	高宏宇
// @param	errno *Errno 基础错误类型	err	error 自定义错误
// @return	*Err 自定义详细错误指针
func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

// @title	DecodeErr
// @description	解析定制的错误。由于 Error 和 Err 都实现了 error 接口中声明的 Error 方法，因此两种类型的实例都可以向上转型为 error 接口类型的变量 err。如果没有接口这种上转型实现的多态，Error和Err类型的解析就需要写2个函数来解决
// @auth	高宏宇
// @param	err	error 自定义错误
// @return	int 业务响应状态码	string 业务响应消息
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	// 断言判断类型的方法只能在switch中使用
	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, err.Error()
}

/* ------------------------------------------------------------------------------ */

// func NewErrno() *Errno {
// 	return &Errno{Code: 1, Message: "error message is null"}
// }

// func IsErrUserNotFound(err error) bool {
// 	code, _ := DecodeErr(err)
// 	return code == ErrUserNotFound.Code
// }
