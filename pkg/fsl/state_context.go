package fsl

// StateContext is the represent of current state's information
type StateContext interface {
	// Input will return the state's raw input data.
	Input() []byte

	// ContextData will return the state's context object data.
	ContextData() []byte

	// EffectiveInput will return the state's effective input, which
	// is the input after InputPath、Parameter have applied.
	EffectiveInput() []byte

	// Output will return the state's raw output data.
	Output() []byte

	// EffectiveOutpu will return the state's effective output, which
	// is the output after ResultSelector、ResultPath、OutputPath have applied.
	EffectiveOutput() []byte
}
