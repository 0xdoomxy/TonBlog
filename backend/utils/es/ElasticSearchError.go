package es

type ESResponseError struct {
}

// Error returns the error message.
// es 查询会返回 resp和err 两个错误，当err为空但是resp的status code不为200时，为了保持系统层面的开发一致性，内置一个ESResponseError
func (e *ESResponseError) Error() string {
	return "es response error"
}
