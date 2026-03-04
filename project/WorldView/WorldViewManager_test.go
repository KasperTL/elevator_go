package WorldView

import (
	"testing"
)

func TestInitWorldView(t *testing.T) {
	myNodeID := 0
	myWorldView := InitWorldView(myNodeID)

	if myWorldView.SenderID != myNodeID {
		t.Errorf("Expected SenderID %d, got %d", myNodeID, myWorldView.SenderID)
	}

	t.Log("Successfully initialized WorldView")
}

func TestSyncOnRejon(t *testing.T) {
	myNodeID := 0
	myWorldView := InitWorldView(myNodeID)
	alivePeers := []int{0, 1, 2}

	// Test syncOnRejon function with empty orders
	orders := myWorldView.Orders
	result := syncOnRejon(orders, alivePeers)

	if result == nil {
		t.Errorf("Expected non-nil orders, got nil")
	}

	t.Log("Successfully tested syncOnRejon")
}

func TestUpdatePeerStatusInMyWorldView(t *testing.T) {
	myNodeID := 0
	myWorldView := InitWorldView(myNodeID)

	peerWorldView := InitWorldView(1)
	result := updatePeerStatusInMyWorldView(myWorldView, peerWorldView)

	if result.SenderID != myNodeID {
		t.Errorf("Expected SenderID %d, got %d", myNodeID, result.SenderID)
	}

	t.Log("Successfully tested updatePeerStatusInMyWorldView")
}
