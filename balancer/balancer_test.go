package balancer

import (
	"testing"
)

func TestQueueBasic(t *testing.T) {
	req := NewRequest(1, 1)
	req2 := NewRequest(2, 1)
	/////// Test basic push, push, pop, pop ///////
	q := make(queue, 0)
	q = q.push(req)
	q = q.push(req2)
	popped, q := q.pop()
	if popped.id != 1 {
		t.Errorf("Failed to pop 'first', popped '%v' instead", popped.id)
	}
	popped, q = q.pop()
	if popped.id != 2 {
		t.Errorf("Failed to pop 'second', popped '%v' instead", popped.id)
	}

	//////// Test basic push pop //////
	q = make(queue, 0)
	q = q.push(NewRequest(3, 1))
	popped, q = q.pop()
	if popped.id != 3 {
		t.Errorf("Expected 'hi' to be popped, received '%v' instead", popped)
	}

	////// Popping empty queue ///////
	q = make(queue, 0)
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}
	q = q.push(NewRequest(4, 1))
	popped, q = q.pop()
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%v'", popped)
	}

}
