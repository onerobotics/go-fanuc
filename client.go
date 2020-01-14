package fanuc

type ErrorClient interface {
	Errors() ([]Error, error)
}

type FrameClient interface {
	Frames() ([]Frame, error)
}

type IOClient interface {
	IO(...Type) ([]IO, error)
}

type NumericRegisterClient interface {
	NumericRegisters() ([]NumericRegister, error)
}

type PositionRegisterClient interface {
	PositionRegisters() ([]PositionRegister, error)
}

type Client interface {
	ErrorClient
	FrameClient
	IOClient
	NumericRegisterClient
	PositionRegisterClient
}
