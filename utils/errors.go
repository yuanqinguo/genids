package utils

// 定义错误码
type Errno struct {
	Code    int64
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

var (
	Ok           = &Errno{Code: 100, Message: "OK"}
	NextIdErr    = &Errno{Code: 101, Message: "Id生成失败"}
	IDInvalidErr = &Errno{Code: -99, Message: "Node ID不合法，只能为0,1,2,且与其他机器不重复"}
	IDWorkErr    = &Errno{Code: -1, Message: "IdWorker对象初始化失败，请重新使用init_node初始化后重试"}
)
