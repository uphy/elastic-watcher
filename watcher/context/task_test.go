package context

import "fmt"

func ExampleTaskSimple() {
	ctx := TODO()
	runner := NewTaskRunner(ctx)
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
	runner := NewTaskRunner(ctx)
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task1")
		ctx.SetPayload([]interface{}{1, 2})
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		fmt.Printf("task2-%v\n", p)
		return nil
	})
	// Output:
	// task1
	// task2-1
	// task2-2
}

func ExampleTaskStop() {
	ctx := TODO()
	runner := NewTaskRunner(ctx)
	runner.RunFunc(func(ctx ExecutionContext) error {
		fmt.Println("task1")
		ctx.SetPayload([]interface{}{1, 2})
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		fmt.Printf("task2-%v\n", p)
		if n, ok := p.(int); ok {
			if n == 2 {
				return ErrStop
			}
		}
		return nil
	})
	runner.RunFunc(func(ctx ExecutionContext) error {
		p := ctx.Payload()
		fmt.Printf("task3-%v\n", p)
		return nil
	})
	// Output:
	// task1
	// task2-1
	// task2-2
	// task3-1
}
