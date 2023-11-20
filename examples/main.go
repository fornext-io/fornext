package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/fornext-io/fornext/pkg/asl"
	"github.com/fornext-io/fornext/pkg/executor"
	"github.com/fornext-io/fornext/test"
)

func main() {
	data, err := os.ReadFile(path.Join(test.CurrentProjectPath(), "examples", "example.json"))
	if err != nil {
		panic(err)
	}

	var sm asl.StateMachine
	err = json.Unmarshal(data, &sm)
	if err != nil {
		panic(err)
	}

	e := executor.NewExecutor(&sm)
	e.Run()
	ec := e.WaitExecutionDone()
	fmt.Printf("id: %v\n", ec.ID)
}
