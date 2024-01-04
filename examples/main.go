// Package main ...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/fornext-io/fornext/pkg/executor"
	"github.com/fornext-io/fornext/pkg/fsl"
	"github.com/fornext-io/fornext/test"
)

func main() {
	data, err := os.ReadFile(path.Join(test.CurrentProjectPath(), "test", "resources", "example.json"))
	if err != nil {
		panic(err)
	}

	var sm fsl.StateMachine
	err = json.Unmarshal(data, &sm)
	if err != nil {
		panic(err)
	}

	e := executor.NewExecutor(&sm)
	e.Run([]byte(`{"name": "123", "data": 10}`))
	ec := e.WaitExecutionDone()
	fmt.Printf("id: %v\n", ec.ID)
}
