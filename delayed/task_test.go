package delayed

import (
	"context"
	"fmt"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestTaskScheduler(t *testing.T) {
	scheduler, err := NewScheduler("root:@/televi?parseTime=true")
	if err != nil {
		panic(err)
	}
	Register(scheduler, "some-task", func(args int) {
		fmt.Println("Called trigger with ", args)
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			fmt.Println("Planning task")
			err = scheduler.Schedule("some-task", t.Add(time.Second), rand.Intn(10))
			fmt.Println(err)
		}
	}()
	scheduler.Run(ctx)
}
