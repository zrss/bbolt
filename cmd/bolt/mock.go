package main

import (
	"fmt"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("frag.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}
		return err
	})*/

	go func() {
		db.View(func(tx *bolt.Tx) error {
			s := db.Stats()
			fmt.Printf("read txn txid: %d. pending: %d, free: %d, open: %d\n", tx.ID(), s.PendingPageN, s.FreePageN, s.OpenTxN)

			fmt.Printf("start of long run read txn\n")
			fmt.Printf("read txn txid: %d\n", tx.ID())
			bucket := tx.Bucket([]byte("MyBucket"))
			bucket.Get([]byte("answer"))

			<-time.After(10 * time.Second)
			fmt.Printf("end of long run read txn\n")
			return nil
		})
	}()

	mockValue := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		time.Sleep(1 * time.Second)
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			err = b.Put([]byte("answer"), mockValue)
			return err
		})
	}

	<-time.After(20 * time.Second)
}
