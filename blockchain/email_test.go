package blockchain

import (
	"bytes"
	"testing"
)

func TestEmail_GenerateID(t *testing.T) {
	// CASE 1: id of email should be unique
	email1 := Email{Sender:"1",Recipient:"2",Message:"hello world"}
	email2 := Email{Sender:"2",Recipient:"1",Message:"hello world"}
	email1.GenerateID()
	email2.GenerateID()
	if email1.ID == email2.ID{
		t.Errorf("ID should not match")
	}
	// CASE 2: id of email should be same if content of email is identical
	email3 := Email{Sender:"1",Recipient:"2",Message:"hello world"}
	email3.GenerateID()
	if email1.ID != email3.ID{
		t.Errorf("ID should  match but got %s & %s",email1.ID,email3.ID)
	}
}