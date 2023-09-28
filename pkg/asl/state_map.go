package asl

import (
	"context"
)

// MapState causes the interpreter to process all the elements of an array,
// potentially in parallel, with the processing of each element independent of the others.
type MapState struct {
	// Type is the type's name of MapState
	// +Required
	// MUST be `Map`
	Type StateType `json:"Type"`

	// Comment provided for human-readable description of the state.
	// +Optional
	Comment string `json:"Comment,omitempty"`

	// InputPath is an Path, which is applied to a State's raw input to select some
	// or all of it, that selection is used by the state.
	// +Optional
	// Defaults to '$'
	InputPath *Path `json:"InputPath,omitempty"`

	// OutputPath is an Path, which is applied to the state's output after the application of `ResultPath`,
	// producing the effective output which serves as the raw input for the next state.
	// +Optional
	// Defaults to '$'
	OutputPath *Path `json:"OutputPath,omitempty"`

	// Next is the name of state the interpreter follows a transition to.
	// It MUST exactly and case-sensitively match the name of the another state.
	// +Optional
	Next string `json:"Next,omitempty"`

	// End causes the interpreter to terminate the machine.
	// +Optional
	End bool `json:"End,omitempty"`

	// ResultPath is a Reference Path, which specifies the raw input's combination with or
	// replacement by the state's result.
	// The value of `ResultPath` MUST NOT begin with '$$'
	// +Optional
	// Defaults to '$'
	ResultPath *ReferencePath `json:"ResultPath,omitempty"`

	// Parameters is a Payload Template which is a JSON object, whose input is the result of
	// applying the `InputPath` to the raw input.
	// If the `Parameters` is provided, its payload, after the extraction and embedding,
	// becomes the effective result.
	// +Optional
	Parameters *PayloadTemplate `json:"Parameters,omitempty"`

	// ResultSelector is a Payload Template, whose input is the result, and whose
	// payload replaces and becomes the effective result.
	// +Optional
	ResultSelector *PayloadTemplate `json:"ResultSelector,omitempty"`

	// Retry is an array of Retrier, when a state reports an error, the interpreter scans through the
	// Retriers and, when the Error Name appears in the value of
	// a Retrier's `ErrorEquals` field, implements the retry policy described in that Retrier.
	// +Optional
	Retry Retriers `json:"Retry,omitempty"`

	// Catch is an array of Catcher, when a state reports an error and either there is no Retrier,
	// or retries have failed to resolve the error, the interpreter scans through the Catchers in
	// array order, and when the Error Name appears in the value of a Catcher's `ErrorEquals` field,
	// transitions the machine to the state named in the value of the `Next` field.
	// +Optional
	Catch Catchers `json:"Catch,omitempty"`

	// ------ State specified field

	// ItemsPath is a Reference Path and MUST identify a field whose value is a JSON array.
	// Defaults to '$'
	// +Optional
	ItemsPath *ReferencePath `json:"ItemsPath,omitempty"`

	// ItemReader is an JSON Object causes the interpreter to read items from the source identified by the
	// ItemReader Configuration's `Resource` field.
	// See more https://docs.aws.amazon.com/step-functions/latest/dg/input-output-itemreader.html
	// +Optional
	// Defaults to '$'
	// DIDN'T Support except Defaults value.
	ItemReader *ItemReader `json:"ItemReader,omitempty"`

	// ItemSelector is a Payload Template, the interpreter uses uses this field to
	// overwide each single element of the item array.
	// Replace `Parameters` field.
	// Defaults to a single element of the array field.
	// +Optional
	ItemSelector *PayloadTemplate `json:"ItemSelector,omitempty"`

	// ItemBatcher is a JSON Object called the ItemBatcher Configuration, cause
	// the interpreter to batch selected items into sub-arrays before passing them to each
	// invocation.
	// See more https://docs.aws.amazon.com/step-functions/latest/dg/input-output-itembatcher.html
	// Defaults to the selected item.
	// +Optional
	ItemBatcher *ItemBatcher `json:"ItemBatcher,omitempty"`

	// ResultWriter is a JSON Object called the ResultWriter Congifuration, cause
	// the interpreter to write results to the resource identified by the `Resource` field.
	// +Optional
	// DIDN'T Support except Defaults value.
	ResultWriter *ItemWriter `json:"ResultWriter,omitempty"`

	// MaxConcurrency is a non-negatie integer value, the interpreter will not allow the number
	// of concurrent iterations to exceed that value.
	// +Optional
	// Defaults to 0
	MaxConcurrency *int `json:"MaxConcurrency,omitempty"`

	// MaxConcurrencyPath MUST be a Reference Path to the MaxConcurrency data.
	// +Optional
	MaxConcurrencyPath *ReferencePath `json:"MaxConcurrencyPath,omitempty"`

	// ToleratedFailurePercentage MUST be a number between zero and 100,
	// the interpreter will continue starting iterations even if some items fail.
	// +Optional
	// Defaults to 0
	ToleratedFailurePercentage *int `json:"ToleratedFailurePercentage,omitempty"`

	// ToleratedFailurePercentagePath MUST be a Reference Path to the ToleratedFailurePercentage data.
	// +Optional
	ToleratedFailurePercentagePath *ReferencePath `json:"ToleratedFailurePercentagePath,omitempty"`

	// ToleratedFailureCount MUST be a non-negative integer.
	// +Optional
	// Defaults to 0
	ToleratedFailureCount *int `json:"ToleratedFailureCount,omitempty"`

	// ToleratedFailureCountPath MUST be a Reference Path to the ToleratedFailureCount data.
	// +Optional
	ToleratedFailureCountPath *ReferencePath `json:"ToleratedFailureCountPath,omitempty"`

	// ItemProcessor specify the `Map` state processing mode and definition.
	// Replace `Iterator` field.
	// +Required
	ItemProcessor ItemProcessor `json:"ItemProcessor"`
}

// ItemProcessor specify the `Map` state processing mode and definition.
type ItemProcessor struct {
	// ProcessorConfig specifies the processing mode for the `Map` state.
	// +Optional
	// DIDN'T Support except Defaults value.
	ProcessorConfig *ProcessorConfig `json:"ProcessorConfig,omitempty"`

	// StartAt is the name of state which inerpreter will start running.
	// The value MUST exactly match one of names of the `States` field.
	// +Required
	StartAt string `json:"StartAt"`

	// States represent the states.
	// +Required
	States States `json:"States"`
}

// ProcessorConfig specifies the processing mode for the `Map` state.
type ProcessorConfig struct {
}

// ItemReader causes the interpreter to read items from the source.
type ItemReader struct {
}

// ItemWriter causes the interpreter write results to the dist.
type ItemWriter struct {
}

// ItemBatcher cause the interpreter to batch selected items.
type ItemBatcher struct {
	// BatchInput is a Payload Template to include in each batch passed to each child workflow execution.
	// The interpreter merges this input with the input for each individual child workflow executions
	// +Optional
	BatchInput *PayloadTemplate `json:"BatchInput,omitempty"`

	// MaxItemsPerBatch MUST be a positive integer, specifies the maximum number of items
	// that each child workflow execution processes.
	// +Optional
	MaxItemsPerBatch *int `json:"MaxItemsPerBatch,omitempty"`

	// MaxInputBytesPerBatch MUST be a positive integer, specifies the maximum size of a batch in bytes.
	// +Optional
	// up to 256 KBs.
	MaxInputBytesPerBatch *int `json:"MaxInputBytesPerBatch,omitempty"`

	// MaxItemsPerBatchPath is a Reference Path to the MaxItemsPerBatch data.
	// +Optional
	MaxItemsPerBatchPath *ReferencePath `json:"MaxItemsPerBatchPath,omitempty"`

	// MaxInputBytesPerBatchPath is a Reference Path to the MaxInputBytesPerBatch data.
	// +Optional
	MaxInputBytesPerBatchPath *ReferencePath `json:"MaxInputBytesPerBatchPath,omitempty"`
}

// Validate will validate the MapState configuration
func (s *MapState) Validate(_ context.Context) error {
	return nil
}
