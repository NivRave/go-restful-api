package user

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/asdine/storm/v3"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func TestCRUD(t *testing.T) {
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Dave",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	t.Log("Read")
	u2, err := GetOne(u.ID)
	if err != nil {
		t.Fatalf("Error reading a record: %s", err)
	}
	if !reflect.DeepEqual(u2, u) {
		t.Fatalf("Records don't match")
	}

	t.Log("Update")
	u.Role = "developer"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	t.Log("Read")
	u3, err := GetOne(u.ID)
	if err != nil {
		t.Fatalf("Error reading a record: %s", err)
	}
	if !reflect.DeepEqual(u3, u) {
		t.Fatalf("Records don't match")
	}

	t.Log("Delete")
	err = DeleteOne(u.ID)
	if err != nil {
		t.Fatalf("Error deleting a record: %s", err)
	}
	_, err = GetOne(u.ID)
	if err == nil {
		t.Fatalf("Record should not exist after deletion")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("%s", "Error retrieving non-existing record: "+err.Error()+"Expected: "+storm.ErrNotFound.Error())
	}
	t.Log("Read All")
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	users, err := GetAll()
	if err != nil {
		t.Fatalf("Error reading all records: %s", err)
	}
	if len(users) != 3 {
		t.Fatalf("Received different number of records. Expected 3, received %d", len(users))
	}
}

func BenchmarkCRUD(b *testing.B) {
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Dave",
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		_, err = GetOne(u.ID)
		if err != nil {
			b.Fatalf("Error reading a record: %s", err)
		}

		u.Role = "developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		err = DeleteOne(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a record: %s", err)
		}
	}
}

func BenchmarkCreate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Dave_" + strconv.Itoa(i),
			Role: "Tester",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Dave_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		_, err = GetOne(u.ID)
		if err != nil {
			b.Fatalf("Error reading a record: %s", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Dave_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		u.Role = "developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Dave_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		err = DeleteOne(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a record: %s", err)
		}
	}
}

func cleanDb(b *testing.B) {
	os.Remove(dbPath)
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Dave",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		b.Fatalf("Error saving a record: %s", err)
	}
	b.ResetTimer()
}
