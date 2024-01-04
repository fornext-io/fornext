package fsl

const (
	// ErrorCodeStatesAll is a wildcard which matches any Error Name
	ErrorCodeStatesAll string = "States.ALL"

	// ErrorCodeStatesHeartbeatTimeout represent a Task State failed to heartbeat for
	// a time longer than the "HeartbeatSeconds" value
	ErrorCodeStatesHeartbeatTimeout string = "States.HeartbeatTimeout"

	// ErrorCodeStatesTimeout represent a Task State either ran longer than the
	// "TimeoutSeconds" value, or failed to heartbeat for a time longer than
	// the "HeartbeatSeconds" value
	ErrorCodeStatesTimeout string = "States.Timeout"

	// ErrorCodeStatesTaskFailed represent a Task State failed during the execution
	ErrorCodeStatesTaskFailed string = "States.TaskFailed"

	// ErrorCodeStatesPermissions represent a Task State failed because it had
	// insufficient privileges to execute the specified code
	ErrorCodeStatesPermissions string = "States.Permissions"

	// ErrorCodeStatesResultPathMatchFailure represent a state's "ResultPath" field
	// cannot be applied to the input the state received
	ErrorCodeStatesResultPathMatchFailure string = "States.ResultPathMatchFailure"

	// ErrorCodeStatesParameterPathFailure represent within a state's "Parameters" field,
	// the attempt to replace a field whose name ends in ".$" using a Path failed
	ErrorCodeStatesParameterPathFailure string = "States.ParameterPathFailure"

	// ErrorCodeStatesBranchFailed represent a branch of a Parallel State failed
	ErrorCodeStatesBranchFailed string = "States.BranchFailed"

	// ErrorCodeStatesNoChoiceMatched represent a Choice State failed to find a match for the
	// condition field extracted from its input
	ErrorCodeStatesNoChoiceMatched string = "States.NoChoiceMatched"

	// ErrorCodeStatesIntrinsicFailure represent within a Payload Template, the attempt
	// to invoke an Intrinsic Function failed
	ErrorCodeStatesIntrinsicFailure string = "States.IntrinsicFailure"

	// ErrorCodeStatesExceedToleratedFailureThreshold represent a Map state failed
	// because the number of failed items exceeded the configured tolerated failure threshold.
	ErrorCodeStatesExceedToleratedFailureThreshold string = "States.ExceedToleratedFailureThreshold"

	// ErrorCodeStatesItemReaderFailed represent a Map state failed to
	// read all items as specified by the `ItemReader` field.
	ErrorCodeStatesItemReaderFailed string = "States.ItemReaderFailed"

	// ErrorCodeStatesResultWriterFailed represent a Map state failed to
	// write all results as specified by the `ResultWriter` field.
	ErrorCodeStatesResultWriterFailed string = "States.ResultWriterFailed"
)

const (
	// ErrorCodeStatesPrefix is the prefix for pre defined ErrorCode
	ErrorCodeStatesPrefix = "States."
)

var (
	// AllErrorCodes contains all pre defined error codes
	AllErrorCodes = map[string]bool{
		ErrorCodeStatesAll:                             true,
		ErrorCodeStatesHeartbeatTimeout:                true,
		ErrorCodeStatesTimeout:                         true,
		ErrorCodeStatesTaskFailed:                      true,
		ErrorCodeStatesPermissions:                     true,
		ErrorCodeStatesResultPathMatchFailure:          true,
		ErrorCodeStatesParameterPathFailure:            true,
		ErrorCodeStatesBranchFailed:                    true,
		ErrorCodeStatesNoChoiceMatched:                 true,
		ErrorCodeStatesIntrinsicFailure:                true,
		ErrorCodeStatesExceedToleratedFailureThreshold: true,
		ErrorCodeStatesItemReaderFailed:                true,
		ErrorCodeStatesResultWriterFailed:              true,
	}
)
