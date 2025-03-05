package transfer

// TransferError 传输错误
type TransferError struct {
	Operation string // 操作名称
	Path      string // 文件路径
	Err       error  // 原始错误
}

// Error 实现error接口
func (e *TransferError) Error() string {
	if e.Path == "" {
		return e.Operation + ": " + e.Err.Error()
	}
	return e.Operation + " " + e.Path + ": " + e.Err.Error()
}

// Unwrap 返回原始错误
func (e *TransferError) Unwrap() error {
	return e.Err
}

// NewTransferError 创建传输错误
func NewTransferError(op, path string, err error) error {
	return &TransferError{
		Operation: op,
		Path:      path,
		Err:       err,
	}
}
