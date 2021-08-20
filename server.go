package main

import (
	"context"
	"encoding/json"
	"login_server/data"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const privateKey = "efjACGRY#WhxARaQ_Fhgm9Vp@zq=kn2Pn8$LNeqFcm#UZ3t7h?Bn@+Z?LsyWYatw"

type LoginServer struct {
	db     *data.MongoUserDB
	router *httprouter.Router
	port   int
}

// NewLoginServer ... Construct a new login server with the given parameters.
func NewLoginServer(port int, mongoURI string) LoginServer {

	db, err := data.NewMongoUserDB(context.Background(), mongoURI)
	if err != nil {
		zap.L().Fatal("Could not create mongo databse.", zap.Error(err))
	}

	ls := LoginServer{
		&db,
		httprouter.New(),
		port,
	}

	ls.setupRoutes()

	return ls
}

// Start ... Start the login server.
func (ls *LoginServer) Start() error {
	return http.ListenAndServe(":"+strconv.Itoa(ls.port), ls.router)
}

func (ls *LoginServer) setupRoutes() {
	ls.router.POST("/auth", ls.authenticate)
	ls.router.POST("/createuser", ls.createUser)
	ls.router.POST("/", ls.defaultRoute)
}

func getSignedKey(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"status":    "OK",
		"user":      userID,
		"ExpiresAt": 15000,
	})

	return token.SignedString([]byte(privateKey))
}

func (ls *LoginServer) defaultRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	zap.L().Warn("Unhandled url", zap.String("url", r.URL.String()))
}

func (ls *LoginServer) authenticate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	select {
	case <-r.Context().Done():
		zap.L().Info("Request handler was cancelled by the timeout handler.")
	default:
		// Metrics.
		start := time.Now()

		// Don't let the client send more than a kilobyte.
		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		// Create a new json decoder.
		dec := json.NewDecoder(r.Body)

		// Don't let the client send fields we aren't expecting.
		dec.DisallowUnknownFields()

		var loginReq data.LoginRequest
		err := dec.Decode(&loginReq)

		// Couldn't decode the message, nothing to do here...
		if err != nil {
			zap.L().Error("Bad request from client", zap.Error(err))
			http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		zap.L().Info("Login Attempt", zap.String("userName", loginReq.Auth.Email))

		// Try to get the user from the database.
		user, err := ls.db.GetUser(r.Context(), loginReq.Auth.Email)

		// If we could not find the user.
		if err != nil {
			zap.L().Info("Login Attempt Failed", zap.String("userName", loginReq.Auth.Email), zap.String("reason", "invalid user name"))
			http.Error(w, "Credentials Rejected", http.StatusForbidden)
			// Nothing else to do here..
			return
		}

		// Now that we have the user we can check if the password is correct.
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Auth.Password))

		// If the password was not correct...
		if err != nil {
			zap.L().Info("Login Attempt Failed", zap.String("userName", loginReq.Auth.Email), zap.String("reason", "invalid password"))
			w.WriteHeader(http.StatusForbidden)
			http.Error(w, "Credentials Rejected", http.StatusForbidden)
			return
		}

		// Generate the token using the secrete key.
		tokenStr, err := getSignedKey(user.ID.String())

		// If that failed panic and let teh globale xception handler catch it.
		if err != nil {
			panic(err)
		}

		// Set the response in the token.
		resp := data.LoginResponse{Token: tokenStr}

		// Write the response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		// Log the completion of the handler.
		zap.L().Info("Login attempt Succedded", zap.Duration("executionTime", time.Since(start)), zap.String("userName", loginReq.Auth.Email))
	}
}

func (ls *LoginServer) createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	select {
	case <-r.Context().Done():
		zap.L().Info("Request handler was cancelled by the timeout handler.")
	default:
		// Metrics.
		start := time.Now()

		// Don't let the client send more than a kilobyte.
		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		// Create a new json decoder.
		dec := json.NewDecoder(r.Body)

		// Don't let the client send fields we aren't expecting.
		dec.DisallowUnknownFields()

		var createUserReq data.CreateUserRequest
		err := dec.Decode(&createUserReq)

		// Couldn't decode the message, nothing to do here...
		if err != nil {
			zap.L().Error("Bad request from client", zap.Error(err))
			http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		zap.L().Info("Create User Request", zap.Any("request", createUserReq))

		// Construct our response.
		createUserResp := data.CreateUserResponse{Success: true, Reason: "OK"}

		// Check the db to see if the user already exists.
		_, err = ls.db.GetUser(r.Context(), createUserReq.Auth.Email)

		// If we found a user we can't actually make a new one so return already exists.
		if err == nil {
			createUserResp.Success = false
			createUserResp.Reason = "User already exists."
			zap.L().Info("Client requested to create a new user that already existed.", zap.Any("request", createUserReq))
		} else {
			// // Okay now we can create the user in teh database.
			userID, err := ls.db.CreateUser(r.Context(), &createUserReq)

			if err != nil {
				panic(err)
			}

			tokenStr, err := getSignedKey(userID.String())

			if err != nil {
				panic(err)
			}

			createUserResp.Token = tokenStr
		}

		zap.L().Info("Create user requests completed", zap.Duration("executionTime", time.Since(start)), zap.Any("response", createUserResp))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createUserResp)
	}
}
