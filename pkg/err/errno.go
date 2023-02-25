package errno

import "fmt"

/* ------------------------------------------------------------------------------ */

/*
	Error结构体用于前端显示，Err结构体用于日志信息
	这里的设计想法是：反馈给前端的错误信息 和 写入日志的错误信息 描述细节等级应当是不同的
	因为反馈给前端的错误信息是用于显示的，必然不能包含太多信息，否则容易遭受攻击
	但是写入日志的信息应当详细，以便维护人员查找错误原因
	所以设计了两个错误，Errno用于基础错误信息，Err在基础错误信息基础上添加了一个额外的错误，仅包含详细错误内容
*/

type Errno struct {
	Code    int
	Message string
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	// Err是接口类型
	Err error
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

// 新建定制的错误
func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

// 解析定制的错误
// 由于 Error 和 Err 都实现了 error 接口中声明的 Error 方法，因此两种类型的实例都可以向上转型为 error 接口类型的变量 err
// 个人理解 Error 和 Err 都属于自定义错误，都应当去实现 error接口中声明的Error方法
// 如果没有接口这种上转型实现的多态，Error和Err类型的解析就需要写2个函数来解决
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	// 这里判断类型的方法只能在switch中使用
	// 这里比较神奇是typed的类型是在发生变化的
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
