package errs

const (
	CommonCodeInit Code = 100000 * (iota + 1)
	AccountCodeInit
	AdminCodeInit
	MatchCodeInit
	OrderCodeInit
	QuoteCodeInit
)

const (
	InternalCode = CommonCodeInit + iota + 1
	RedisErrCode
	ExecSqlFailedCode
	ParamValidateFailedCode
	RecordNotFoundErrCode
	DuplicateDataErrCode
	MongoErrCode
	KafkaErrCode
	EtcdErrCode
	DtmErrCode
	PulsarErrCode
	TimoutOutCode
)

const (
	ParamValidateFailed Code = 170000
)

var (
	// 通用错误
	InternalErr            = InternalCode.error()
	RedisErr               = RedisErrCode.error()
	RecordNotFoundErr      = RecordNotFoundErrCode.error()
	TimeOutErr             = TimoutOutCode.error()
	ParamValidateFailedErr = ParamValidateFailed.error()
)
