package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgraph-io/dgo/v240"
)

const (
	qlGetFriendsOneHop = `query friends($a: string) {
		friends(func: uid($a)) {
		  name
		  type
		  created_at
		  friend_of {
			name
		  }
		  ~friend_of {
			name
		  }
		  xid
		}
	  }`

	qlGetFriendsTwoHop = `query friends($a: string) {
		friends(func: uid($a)) {
		  name
		  type
		  created_at
		  friend_of {
			name
			friend_of {
				name
			}
		  }
		  ~friend_of {
			name
			~friend_of {
				name
			}
		  }
		  xid
		}
	  }`
)

func init() {
	BenchTasks["get-friends-one-hop"] = GetFriendsOneHop
	BenchTasks["get-friends-two-hop"] = GetFriendsTwoHop
}

func GetFriendsOneHop(dgraphCli *dgo.Dgraph, r *rand.Rand) error {
	start := time.Now()

	uid := r.Int63n(MaxUid)
	txn := dgraphCli.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(context.Background(), qlGetFriendsOneHop, map[string]string{"$a": strconv.FormatInt(uid, 10)})
	if err != nil {
		counters.WithLabelValues("get-friends-one-hop", "ERROR").Inc()
		fmt.Printf("Failed to query friends for uid %d: %v\n", uid, err)
		return err
	}

	// fmt.Println(string(resp.Json))

	_ = resp.Json

	counters.WithLabelValues("get-friends-one-hop", "OK").Inc()
	durations.WithLabelValues("get-friends-one-hop", "OK").Observe(time.Since(start).Seconds())

	return nil
}

func GetFriendsTwoHop(dgraphCli *dgo.Dgraph, r *rand.Rand) error {
	start := time.Now()

	uid := r.Int63n(MaxUid)
	txn := dgraphCli.NewReadOnlyTxn()
	resp, err := txn.QueryWithVars(context.Background(), qlGetFriendsTwoHop, map[string]string{"$a": strconv.FormatInt(uid, 10)})
	if err != nil {
		counters.WithLabelValues("get-friends-two-hop", "ERROR").Inc()
		fmt.Printf("Failed to query friends for uid %d: %v\n", uid, err)
		return err
	}

	// fmt.Println(string(resp.Json))

	_ = resp.Json

	counters.WithLabelValues("get-friends-two-hop", "OK").Inc()
	durations.WithLabelValues("get-friends-two-hop", "OK").Observe(time.Since(start).Seconds())

	return nil
}
