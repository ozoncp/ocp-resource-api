package saver

import (
	"fmt"
	"github.com/ozoncp/ocp-resource-api/internal/flusher"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"sync"
	"time"
)

type Saver interface {
	Save(entity models.Resource) // заменить на свою сущность
	Close()
	Init()
}

type saver struct {
	Saver
	entities    []models.Resource
	flusher     flusher.Flusher
	flushPeriod time.Duration
	saveChan    chan models.Resource
	closeChan   chan struct{}
	wait        *sync.WaitGroup
	closed      bool
	initiated   bool
}

func (s *saver) Save(entity models.Resource) {
	assertIsReady(s)
	s.saveChan <- entity
}

func (s *saver) Close() {
	assertIsReady(s)
	s.closed = true
	close(s.closeChan)
	s.wait.Wait()
}

func assertIsReady(s *saver) {
	if !s.initiated {
		panic("saver must be initiated")
	}
	if s.closed {
		panic("saver must not be closed")
	}
}

func (s *saver) Init() {
	if s.initiated {
		return
	}
	s.initiated = true
	s.wait.Add(1)
	go func() {
		ticker := time.NewTicker(s.flushPeriod)
		defer ticker.Stop()
		defer s.wait.Done()
		for {
			select {
			case entity := <-s.saveChan:
				s.entities = append(s.entities, entity)
			case <-ticker.C:
				s.flush()
			case <-s.closeChan:
				s.flush()
				if len(s.entities) != 0 {
					fmt.Println("Failed to gracefully close saver")
				}
				return
			}
		}
	}()
}

func (s *saver) flush() {
	fmt.Printf("Flush entities %v\n", s.entities)
	notFlushedEntities := s.flusher.Flush(nil, s.entities, nil)
	if len(notFlushedEntities) != 0 {
		s.entities = notFlushedEntities
		fmt.Printf(
			"Flush failed. %d/%d elements flushed. Remains: %v\n",
			len(s.entities)-len(notFlushedEntities), len(s.entities), notFlushedEntities)
		return
	}
	s.entities = s.entities[:0]
}

func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	flushTimer time.Duration,
) Saver {
	s := saver{
		entities:    make([]models.Resource, 0, capacity),
		flusher:     flusher,
		flushPeriod: flushTimer,
		saveChan:    make(chan models.Resource),
		closeChan:   make(chan struct{}),
		wait:        &sync.WaitGroup{},
	}
	return &s
}
