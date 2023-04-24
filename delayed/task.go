package delayed

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/televi-go/televi/util"
	"log"
	"sync"
	"time"
)

type TaskExecutor[TArgs any] func(args TArgs)

type CommonTaskExecutor func(args []byte)

type ScheduledEntry struct {
	Descriptor string
	ExecuteAt  time.Time
	Args       []byte
	Id         int
}

type TaskScheduler struct {
	Executors map[string][]CommonTaskExecutor
	Scheduled *util.LinkedList[ScheduledEntry]
	mut       sync.RWMutex
	Storage   *sql.DB
}

func NewScheduler(dsn string) (*TaskScheduler, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS TeleviTasks(Descriptor VARCHAR(255) NOT NULL, ExecuteAt TIMESTAMP NOT NULL, Args BLOB, Id int not null primary key AUTO_INCREMENT)")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT Id, Descriptor, ExecuteAt, Args FROM TeleviTasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	scheduler := &TaskScheduler{
		Executors: map[string][]CommonTaskExecutor{},
		Scheduled: util.NewLinkedList[ScheduledEntry](),
		mut:       sync.RWMutex{},
		Storage:   db,
	}

	for rows.Next() {
		var (
			id         int
			descriptor string
			executeAt  time.Time
			argBytes   []byte
		)
		err = rows.Scan(&id, &descriptor, &executeAt, &argBytes)
		if err != nil {
			return nil, err
		}

		scheduler.Scheduled.PushFront(ScheduledEntry{
			Descriptor: descriptor,
			ExecuteAt:  executeAt,
			Args:       argBytes,
			Id:         id,
		})
	}

	return scheduler, nil
}

func (scheduler *TaskScheduler) removeSchedule(descriptor string, id int) {
	scheduler.Storage.Exec("DELETE FROM TeleviTasks WHERE Descriptor = ? AND Id=?", descriptor, id)
}

func Register[T any](scheduler *TaskScheduler, descriptor string, executor TaskExecutor[T]) {
	scheduler.register(descriptor, func(argBytes []byte) {
		var arguments T
		err := json.Unmarshal(argBytes, &arguments)
		if err != nil {
			log.Printf("Error in parsing arguments for %s:%v\n", descriptor, err)
			return
		}
		executor(arguments)
	})
}

func (scheduler *TaskScheduler) register(descriptor string, executor CommonTaskExecutor) {
	scheduler.mut.Lock()
	defer scheduler.mut.Unlock()
	executors, _ := scheduler.Executors[descriptor]
	scheduler.Executors[descriptor] = append(executors, executor)
}

func (scheduler *TaskScheduler) Schedule(descriptor string, executeAt time.Time, args any) error {
	scheduler.mut.Lock()
	defer scheduler.mut.Unlock()

	argsBytes, err := json.Marshal(args)
	if err != nil {
		return err
	}

	result, err := scheduler.Storage.Exec("INSERT INTO TeleviTasks(Descriptor,ExecuteAt, Args) VALUES (?,?,?)", descriptor, executeAt, argsBytes)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	scheduler.Scheduled.PushFront(ScheduledEntry{
		Descriptor: descriptor,
		ExecuteAt:  executeAt,
		Args:       argsBytes,
		Id:         int(id),
	})
	return nil
}

func (scheduler *TaskScheduler) traverseAndRun() {
	scheduler.mut.Lock()
	defer scheduler.mut.Unlock()
	for i, entry := range scheduler.Scheduled.ToArray() {
		if entry.ExecuteAt.After(time.Now()) {
			continue
		}
		scheduler.Scheduled.Delete(i)
		scheduler.removeSchedule(entry.Descriptor, entry.Id)
		executors := scheduler.Executors[entry.Descriptor]
		for _, executor := range executors {
			go executor(entry.Args)
		}
	}
}

func (scheduler *TaskScheduler) Run(ctx context.Context) {
	for {
		select {
		case _, _ = <-ctx.Done():
			return
		default:
			scheduler.traverseAndRun()
			<-time.After(time.Second)
		}
	}
}
