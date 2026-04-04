package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"multiple-copies-go/resourcelock"
)

func main() {
	id := flag.String("id", "", "replica ID")
	lockfile := flag.String("lockfile", "./resourcelock/lock", "path to lock file")
	interval := flag.Duration("interval", 2*time.Second, "heartbeat / retry interval")
	leaseDuration := flag.Duration("lease-duration", 5*time.Second, "lease expiration duration")
	flag.Parse()

	if *id == "" {
		fmt.Println("Error: -id is required")
		os.Exit(1)
	}

	lock := resourcelock.New(*lockfile, *id)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	isLeader := false

	go func() {
		ticker := time.NewTicker(*interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if !isLeader {
					acquired, err := lock.TryLock()
					if err != nil {
						fmt.Printf("[%s] Error acquiring lock: %v\n", *id, err)
						continue
					}
					if acquired {
						isLeader = true
						fmt.Printf("[%s] Elected as LEADER\n", *id)
					} else {
						lease := resourcelock.GetLease(*lockfile)
						if lease == nil {
							fmt.Printf("[%s] Waiting... no holder info\n", *id)
							continue
						}

						// Check if leader's lease has expired
						if !lease.UpdatedAt.IsZero() && time.Since(lease.UpdatedAt) > *leaseDuration {
							fmt.Printf("[%s] Leader '%s' lease expired (last seen: %v ago), grabbing lock!\n",
								*id, lease.HolderID, time.Since(lease.UpdatedAt))
							// Force grab: TryLock will succeed since flock is advisory
							// and the dead process no longer holds the kernel lock.
							acquired, err := lock.TryLock()
							if err != nil {
								fmt.Printf("[%s] Error force-grabbing lock: %v\n", *id, err)
								continue
							}
							if acquired {
								isLeader = true
								fmt.Printf("[%s] Elected as LEADER (after leader expiry)\n", *id)
							}
						} else {
							remaining := *leaseDuration - time.Since(lease.UpdatedAt)
							fmt.Printf("[%s] Waiting... leader: %s, lease expires in %v\n", *id, lease.HolderID, remaining)
						}
					}
				} else {
					if err := lock.Heartbeat(); err != nil {
						fmt.Printf("[%s] Heartbeat failed: %v\n", *id, err)
					}
				}
			}
		}
	}()

	<-sigCh
	fmt.Printf("\n[%s] Shutting down...\n", *id)
	close(done)
	if isLeader {
		// 以 leader 身份正常退出，释放锁
		lock.Unlock()
		fmt.Printf("[%s] Released lock\n", *id)
	}
}
