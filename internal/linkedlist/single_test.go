package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingleListOperations(t *testing.T) {
	s := NewSingle[int]()
	s.Append(1)
	i, ok := s.Shift()
	require.True(t, ok)
	require.Equal(t, i, 1)
	i, ok = s.Shift()
	require.False(t, ok)

	s.Append(1, 2)
	i, ok = s.Shift()
	require.True(t, ok && i == 1)
	i, ok = s.Shift()
	require.True(t, ok && i == 2)
	i, ok = s.Shift()
	require.False(t, ok)

	if s.Len() != 0 {
		t.Errorf("Expected length of 0 but got %d", s.Len())
	}

	s.Append(1, 2, 3)

	if s.Len() != 3 {
		t.Errorf("Expected length of 3 but got %d", s.Len())
	}

	val, ok := s.Shift()

	if !ok {
		t.Errorf("Expected Shift operation to return true but got false")
	}

	if val != 1 {
		t.Errorf("Expected Shift operation to return 1 but got %v", val)
	}

	if s.Len() != 2 {
		t.Errorf("Expected length of 2 but got %d", s.Len())
	}

	s.Append(4)

	if s.Len() != 3 {
		t.Errorf("Expected length of 3 but got %d", s.Len())
	}

	val, ok = s.Shift()

	if !ok {
		t.Errorf("Expected Shift operation to return true but got false")
	}

	if val != 2 {
		t.Errorf("Expected Shift operation to return 2 but got %v", val)
	}

	if s.Len() != 2 {
		t.Errorf("Expected length of 2 but got %d", s.Len())
	}
}
