package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"context"

	"time"

	"github.com/ynishi/indeks"
	"google.golang.org/api/option"
)

func main() {
	DefaultFirebaseClient, _ = MakeClient(context.Background(), "../firebase-admin-sdk.json")

	http.HandleFunc("/check", HealthCheckHandler)
	http.HandleFunc("/idx", AddIdxHandler())
	http.ListenAndServe(":3000", nil)
}

var HealthCheckHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	DefaultFirebaseClient.FirestoreClient.Collection("check")
	fmt.Fprintf(w, "OK")
})

type Client struct {
	Cred            string
	Ctx             context.Context
	App             *firebase.App
	FirestoreClient *firestore.Client
}

var DefaultFirebaseClient = &Client{"", nil, nil, nil}

func MakeClient(ctx context.Context, cred string) (client *Client, err error) {
	client = &Client{}
	opt := option.WithCredentialsFile(cred)
	var mCtx context.Context
	if ctx == nil {
		mCtx = context.Background()
	} else {
		mCtx = ctx
	}

	app, err := firebase.NewApp(mCtx, nil, opt)
	if err != nil {
		return nil, err
	}

	firestoreClient, err := app.Firestore(mCtx)
	if err != nil {
		return nil, err
	}
	client = &Client{cred, mCtx, app, firestoreClient}
	return client, nil
}

func AddIdxHandler() (hf http.HandlerFunc) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("name")
		desc := r.Header.Get("desc")
		//defaultPoint := r.Header.Get("defaultPoint")
		//defaultDuration := r.Header.Get("defaultDuration")

		indeks.Idxs = append(indeks.Idxs, &indeks.Idx{name, desc, 1, time.Duration(1), nil})
		DefaultFirebaseClient.FirestoreClient.Collection("idxs").Add(DefaultFirebaseClient.Ctx, map[string]interface{}{})
		fmt.Fprintf(w, "OK")
	})
}
