package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"restful_api/user"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestBodyToUser(t *testing.T) {
	validUser := &user.User{
		ID:   bson.NewObjectId(),
		Name: "Chad",
		Role: "Tester",
	}

	validUser2 := &user.User{
		ID:   validUser.ID,
		Name: "Chad",
		Role: "Developer",
	}

	json, err := json.Marshal(validUser)
	if err != nil {
		t.Errorf("Error marshalling a valid user:: %s", err)
		t.FailNow()
	}

	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "empty request body",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "empty user",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		},
		{
			txt: "malformed data",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{"id":12}`)),
			},
			u:   &user.User{},
			err: true,
		},
		{
			txt: "valid request",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(json)),
			},
			u:   &user.User{},
			exp: validUser,
		},
		{
			txt: "valid partial request",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{"role":"Developer"}`)),
			},
			u:   validUser,
			exp: validUser2,
		},
	}

	for _, tc := range ts {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expected error, got none")
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different than expected:")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
