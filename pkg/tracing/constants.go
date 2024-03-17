package tracing

const (
	BindType       = "error.bind"
	QueryType      = "error.query"
	FileType       = "error.file"
	DecoderType    = "error.decoder"
	InternalErr    = "error.internal"
	AccessType     = "error.access"
	SessionType    = "error.JWT"
	ValidationType = "error.validation"
	TimeFormatType = "error.query-time-format"

	// Managers
	GetFulCaseByIDType      = "error.get-ful-case-by-id"
	GetSpecialistRatingType = "error.get-specialist-rating"

	// Public
	ManagerLoginType       = "error.manager-login"
	SpecialistRegisterType = "error.specialist-register"
	SpecialistLoginType    = "error.specialist-login"
	CameraCreateType       = "error.camera-create"
	CameraDeleteType       = "error.camera-delete"
	CaseCreateType         = "error.case-create"
	RefreshType            = "error.refresh"

	// Specialists
	CreateRatedType     = "error.rated-create"
	GetCasesByLevelType = "error.get-cases"
	GetRatingType       = "error.get-rating"
	GetRatedSolvedType  = "error.get-rated-solved"
	GetMeType           = "error.get-me"
	UpdateMeType        = "error.update-me"
)

const (
	// Managers
	GetFulCaseByID      = "Get ful case info by it's id"
	GetSpecialistRating = "Get specialist rating"

	// Public
	ManagerLogin       = "Manager login"
	SpecialistLogin    = "Specialist login"
	SpecialistRegister = "Specialist register"
	CameraCreate       = "Camera create"
	CameraDelete       = "Camera delete"
	CaseCreate         = "Case create"
	Refresh            = "Refresh"

	// Specialists
	CreateRated     = "Create rated"
	GetCasesByLevel = "Get cases by level"
	GetRating       = "Get rating"
	GetRatedSolved  = "Get rated solved"
	GetMe           = "Get me"
	UpdateMe        = "Update me"
)

const (
	SuccessfulCompleting = "Operation completed successfully"
)
