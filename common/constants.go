package common

var (
	ABHIDomain  = ""
	ABHIAppName = ""
)

const (
	//--------------WALL APPLICATION CONSTANTS ------------------------

	//DEV
	// ABHIDomain = "localhost"
	// ABHIAllowOrigin = "http://localhost:8080"
	// ABHIAppName = "novodev"

	// ABHIDomain = "flattrade.in"
	// ABHIAllowOrigin = "https://novo.flattrade.in"
	// ABHIAppName = "novo"

	ABHICookieName       = "ftab_pt"
	ABHIClientCookieName = "ftab_ud"
	//--------------OTHER COMMON CONSTANTS ------------------------
	CookieMaxAge = 300

	TechExcelPrefix = "TECHEXCELPROD.capsfo.dbo."

	SuccessCode  = "S" //success
	ErrorCode    = "E" //error
	LoginFailure = "I" //??

	StatusPending = "P" //pending
	StatusApprove = "A" //Approve
	StatusReject  = "R" //Reject
	StatusNew     = "N" //new
	Statement     = "1"
	Detail        = "2"
	Panic         = "P"
	NoPanic       = "NP"
	INSERT        = "INSERT"
	UPDATE        = "UPDATE"
	SUCCESS       = "success"
	FAILED        = "failed"
	PENDING       = "pending"
	NSE           = "NSE"
	BSE           = "BSE"
	AUTOBOT       = "AUTOBOT"
	//added by naveen
	Mobile = "M"
	Web    = "W"
)

var ABHIAllowOrigin []string

// var ABHIBrokerId = 0

var ABHIBrokerId = 0

var ABHIFlag string
