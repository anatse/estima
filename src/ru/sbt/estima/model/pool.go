package model

import (
	"fmt"
	"sync"
)

type poolElement struct {
	inUse bool
	object interface{}
}

type CRF func ()(interface{}, error)
type RRF func (interface{})(error)
type ComputePF func (interface{})

type Pool struct {
	sync.Mutex
	maxSize uint
	currentSize uint
	objects []poolElement

	createResource CRF
	releaseResource RRF
}

func NewPool (size uint, crf CRF, rrf RRF) *Pool {
	pool := Pool {
		maxSize: size,
		currentSize: 0,
		createResource: crf,
		releaseResource: rrf,
		objects: make([]poolElement, size),
	}

	return &pool
}

func (pool *Pool) Use (cpf ComputePF) (err error) {
	//log.Printf("Allocating pool for: %v, current size: %d\n", cpf, pool.currentSize)

	pobj, err := pool.Get()
	if err != nil {
		return err
	}
	defer (func() {
		err = pool.Release(pobj)
		//log.Printf("Processed. Released pool for: %v, current size: %d\n", cpf, pool.currentSize)
	})()

	//log.Printf("Allocated pool for: %v, current size: %d, processing...\n", cpf, pool.currentSize)
	cpf (pobj)
	//log.Printf("Processed call. Pool for: %v, current size: %d\n", cpf, pool.currentSize)
	return err
}

func (pool *Pool) Close () {
	pool.Lock()
	defer pool.Unlock()

	for _, pe := range pool.objects {
		pe.inUse = false
		pool.releaseResource(pe.object)
		pool.currentSize--
	}
}

func (pool *Pool) Get()(obj interface{}, err error) {
	pool.Lock()
	defer pool.Unlock()

	if pool.currentSize == pool.maxSize {
		return nil, fmt.Errorf("Pool is full: %d of %d", pool.currentSize, pool.maxSize)
	}

	for idx, pe := range pool.objects {
		if !pe.inUse {
			// Check if object already created
			if pe.object == nil {
				// Create object if empty
				pe.object, err = pool.createResource()
				if err != nil {
					return nil, err
				}
			}

			pe.inUse = true
			pool.currentSize++
			pool.objects[idx] = pe
			return pool.objects[idx].object, nil
		}
	}

	return nil, fmt.Errorf("Unexpected error, Possible cause is wrong written Get function")
}

func (pool *Pool) Release (obj interface{})(err error) {
	pool.Lock()
	defer pool.Unlock()

	for idx, pe := range pool.objects {
		if pe.inUse && pe.object == obj {
			pe.inUse = false
			pool.objects[idx] = pe
			pool.currentSize--
			return nil
		}
	}

	return fmt.Errorf("Pool release error: Object not found")
}
