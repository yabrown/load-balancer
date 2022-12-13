package balancer

import (
	"testing"
)

func TestQueueBasic(t *testing.T) {
	req := newRequest(1)
	req2 := newRequest(2)
	/////// Test basic push, push, pop, pop ///////
	q := make(queue, 0)
	q = q.push(req)
	q = q.push(req2)
	popped, q := q.pop()
	if popped.id != 1 {
		t.Errorf("Failed to pop 'first', popped '%s' instead", popped.id)
	}
	popped, q = q.pop()
	if popped.id != 2 {
		t.Errorf("Failed to pop 'second', popped '%s' instead", popped.id)
	}

	//////// Test basic push pop //////
	q = make(queue, 0)
	q = q.push(newRequest(3))
	popped, q = q.pop()
	if popped.id != 3 {
		t.Errorf("Expected 'hi' to be popped, received '%s' instead", popped)
	}

	////// Popping empty queue ///////
	q = make(queue, 0)
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%s'", popped)
	}
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%s'", popped)
	}
	q = q.push(newRequest(4))
	popped, q = q.pop()
	popped, q = q.pop()
	if popped != nil {
		t.Errorf("Expected 'failed', received '%s'", popped)
	}

}
