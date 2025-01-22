package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
)

const (
	MaxUid = 3300000
)

var (
	personXid int64
)

func init() {
	BenchTasks["insert-friend"] = InsertFriend
	BenchTasks["insert-person"] = InsertPerson
}

func InsertFriend(dgraphCli *dgo.Dgraph, r *rand.Rand) error {
	start := time.Now()

	auid := r.Int63n(MaxUid)
	buid := r.Int63n(MaxUid)
	for auid == buid {
		buid = r.Int63n(MaxUid)
	}

	// fmt.Printf("%d is friend of %d\n", auid, buid)

	person := &Person{
		Uid: strconv.FormatInt(auid, 10),
		FriendOf: &Person{
			Uid: strconv.FormatInt(buid, 10),
		},
	}

	// Insert friend edge
	payload, err := json.Marshal(person)
	if err != nil {
		counters.WithLabelValues("insert-friend", "ERROR").Inc()
		fmt.Printf("Failed to marshal person object: %v\n", err)
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   payload,
	}

	txn := dgraphCli.NewTxn()
	as, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		counters.WithLabelValues("insert-friend", "ERROR").Inc()
		return err
	}

	counters.WithLabelValues("insert-friend", "OK").Inc()
	durations.WithLabelValues("insert-friend", "OK").Observe(time.Since(start).Seconds())
	_ = as

	return nil
}

func InsertPerson(dgraphCli *dgo.Dgraph, r *rand.Rand) error {
	start := time.Now()

	xid := strconv.FormatInt(atomic.AddInt64(&personXid, 1), 10)

	person := &Person{
		Uid:       "_:" + xid,
		Xid:       xid,
		Name:      RandString(10, r),
		CreatedAt: start.Unix(),
		UpdatedAt: start.Unix(),
	}

	// Insert person node
	payload, err := json.Marshal(person)
	if err != nil {
		counters.WithLabelValues("insert-person", "ERROR").Inc()
		fmt.Printf("Failed to marshal person object: %v\n", err)
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   payload,
	}

	txn := dgraphCli.NewTxn()
	as, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		counters.WithLabelValues("insert-person", "ERROR").Inc()
		return err
	}

	counters.WithLabelValues("insert-person", "OK").Inc()
	durations.WithLabelValues("insert-person", "OK").Observe(time.Since(start).Seconds())
	_ = as

	return nil
}
