package context

import "fmt"

func ExampleTaskSimple() {
	ctx := TODO()
	runner := ctx.TaskRunner()
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task1")
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task2")
		return nil
	})
	// Output:
	// task1
	// task2
}

func ExampleTaskSplit() {
	ctx := TODO()
	runner := ctx.TaskRunner()
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task1")
		setSplittedPayload(ctx, []JSONObject{
			{
				"a": 1,
			},
			{
				"a": 2,
			},
		})
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		fmt.Printf("task2-%v\n", p["a"])
		return nil
	})
	// Output:
	// task1
	// task2-1
	// task2-2
}

func ExampleTaskStop() {
	ctx := TODO()
	runner := ctx.TaskRunner()
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task1")
		setSplittedPayload(ctx, []JSONObject{
			{
				"a": 1,
			},
			{
				"a": 2,
			},
		})
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		if fmt.Sprint(p["a"]) == "2" {
			return ErrStop
		}
		fmt.Printf("task2-%v\n", p["a"])
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		fmt.Printf("task3-%v\n", p["a"])
		return nil
	})
	// Output:
	// task1
	// task2-1
	// task3-1
}
