package executor

// ExecutionStartedEvent is result of `StartExecutionCommand`, and MUST followed it.
// When follower's receive this event, it persistent the result to
// storage to keep consistency with leader.
// And if an new leader is elected, and found `StartExecutionCommand` didn't followed by
// any `ExecutionStartedEvent`, treat it as an not processed command and reprocess it.
type ExecutionStartedEvent struct {
	// ID is the identify of this execution.
	ID string

	// WorkflowID is the identify of workflow definition which this execution belongs.
	WorkflowID string

	// Input is the arguments of this execution. It MUST be an valid JSON-Object.
	Input []byte

	// Timestamp is `HLC` value of this command request time.
	Timestamp uint64

	// Status is current state of execution, which could be one of
	// 1. `Pending`: is waiting to start the first `State`, current it always
	//   transition to `Running` right now; but when we support `delay execution`
	//   there will be an observed `Pending` state.
	// 2. `Running`: is running with an state、task and so on.
	// 3. `Completing`: had reach an `End` state, and all running state、task is reclaiming.
	// 4. `Succeeded`: execution is done as `Succeed`.
	// 5. `Failed`: execution is done as `Fail`.
	// 6. `Deleting`: execution record is `Deleting`, all persistent data will been deleted.
	// 7. `Canceling`: user commant to stopping exeuciton, or it's timeout reached.
	// 8. `Canceled`: execution is done as `Cancel`.
	//
	// State transition flow:
	//                                 +-------------+
	//                                 |   Pending   |
	//                                 +------|------+
	//                                        |
	//                                        |
	//                                        |
	//            +-------------+      +------|------+       +-------------+
	//        +----  Completing --------   Running   ---------  Canceling  |
	//        |   +------|------+      +-------------+       +------|------+
	//        |          |                                          |
	//        |          |                                          |
	//        |          |                                          |
	// +------|------+   |  +------+------+                  +------|------+
	// |  Succeeded  |   +---   Failed    ----+              |   Canceled  |
	// +------|------+      +-------------+   |              +------|------+
	//        |                               |                     |
	//        |                               |                     |
	//        |                               |                     |
	//        |                        +------|------+              |
	//        +-------------------------   Deleting  ---------------+
	//                                 +-------------+
	Status string
}
