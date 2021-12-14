package filequeue_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/rstudio/filequeue"
	"github.com/stretchr/testify/require"
)

func TestFileQueue(t *testing.T) {
	r := require.New(t)

	d := t.TempDir()

	q, err := filequeue.New(d)
	r.Nil(err)
	r.NotNil(q)

	qLen, err := q.Len()
	r.Nil(err)
	r.Equal(0, qLen)

	err = q.Push([]byte("kris"))
	r.Nil(err)

	qLen, err = q.Len()
	r.Equal(1, qLen)

	el, err := q.Pop()
	r.Nil(err)
	r.Equal("kris", string(el))

	qLen, err = q.Len()
	r.Equal(0, qLen)

	wg := &sync.WaitGroup{}

	for i := 100; i > 0; i-- {
		for _, item := range []string{"ralsei", "lancer", "flowey"} {
			r.Nil(q.Push([]byte(fmt.Sprintf("%s-%02d", item, i))))
		}
	}

	qLen, err = q.Len()
	r.Equal(300, qLen)

	consumer := func() {
		defer wg.Done()

		name := fmt.Sprintf("maus-%v", rand.Float64())

		q, err := filequeue.New(d)
		r.Nil(err)
		r.NotNil(q)

		for {
			el, err := q.Pop()
			if err != nil {
				t.Logf("consumer=%q pop-missed-err=%v", name, err)
				continue
			}

			if el == nil {
				return
			}

			t.Logf("consumer=%q el=%q", name, string(el))
			time.Sleep(time.Duration(rand.Float64()*10) * time.Millisecond)
		}
	}

	wg.Add(1)
	go consumer()

	wg.Add(1)
	go consumer()

	wg.Add(1)
	go consumer()

	wg.Wait()

	qLen, err = q.Len()
	r.Equal(0, qLen)
}
