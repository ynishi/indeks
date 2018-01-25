package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func setup() {
	DefaultFirebaseClient, _ = MakeClient(context.Background(), "../firebase-admin-sdk.json")
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	os.Exit(ret)
}

func TestInitFirebase(t *testing.T) {
	// FireStore
	ctx := context.Background()

	opt := option.WithCredentialsFile("../firebase-admin-sdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		t.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Fatalf("Firestore client init failed: %v", err)
	}
	defer client.Close()
	_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err != nil {
		t.Fatalf("Failed adding user Ada: %v", err)
	}

}

func TestMakeClient(t *testing.T) {
	client, _ := MakeClient(context.Background(), "../firebase-admin-sdk.json")
	_, _, _ = client.FirestoreClient.Collection("users").Add(client.Ctx, map[string]interface{}{
		"first":  "Alan",
		"middle": "Mathison",
		"last":   "Turing",
		"born":   1912,
	})
}

func TestHealthCheckHttpServer(t *testing.T) {

	ts := httptest.NewServer(HealthCheckHandler)
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

func TestAddIdxAPI(t *testing.T) {

	ts := httptest.NewServer(AddIdxHandler())
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

func TestUpdateIdxAPI(t *testing.T) {

	ts := httptest.NewServer(UpdateIdxHandler())
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

func TestRemoveIdxAPI(t *testing.T) {
	ts := httptest.NewServer(RemoveIdxHandler())
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

func TestListIdxAPI(t *testing.T) {
	ts := httptest.NewServer(ListIdxHandler())
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

func TestDetailIdxAPI(t *testing.T) {
	ts := httptest.NewServer(DetailIdxHandler())
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if "OK" != string(data) {
		t.Fatalf("Data is not OK: %v", string(data))
	}
}

