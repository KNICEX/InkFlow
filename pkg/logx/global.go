package logx

var logger Logger = NewNopLogger()

// Register 注册全局日志, 无锁, 请只在初始化时调用
func Register(l Logger) {
	logger = l
}

func L() Logger {
	return logger
}
