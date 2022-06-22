package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

// ClientHandler ethereum client instance
type Database struct {
	*sql.DB
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateToken(loginJWTToken string) (claims jwt.MapClaims, valid bool) {

	var mySigningKey = []byte("secretkey")
	token, err := jwt.Parse(loginJWTToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Errorf("Your token has expired: %s", err.Error())
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func ClientAPIHandler (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]

	// Get the query parameters from url request
	// address := r.URL.Query().Get("address")
	// hash := r.URL.Query().Get("hash")
	// assetIDs := r.URL.Query().Get("assetIDs")
	// tokenURI := r.URL.Query().Get("tokenURI")

	// Set our  Response header
	w.Header().Set("Content-Type", "application/json")

	// Handle each request using the module parameter:
	switch module {
	case "sign-up":
		var signUpUser SignUpUser
		var dbUser DBUser
		var err error

		err = json.NewDecoder(r.Body).Decode(&signUpUser)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}

		isEmailValid := ValidateEmail(signUpUser.Email)
		if !isEmailValid {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Email not in the correct format",
			})
			return
		}

		row := db.QueryRow("SELECT * FROM users where email= $1", signUpUser.Email)
		err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Passhash, &dbUser.Role)

		if err != nil && err != sql.ErrNoRows {
			fmt.Println("Error querying DB for users", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error querying DB for users",
			})
			return
		}

		//checks if email is already registered
		if dbUser.Email != "" {
			fmt.Println("Email already in use")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Email already in use",
			})
			return
		}

		var newUser = new(DBUser)
		//create newUser details
		newUser.Passhash, err = GeneratehashPassword(signUpUser.Password)
		if err != nil {
			log.Fatalln("error in password hash creation")
		}

		//register new user
		fmt.Println("New User Registration")
		newUser.ID, err = nanoid.New()
		if err != nil {
			fmt.Println("Error generating UUID", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error generating UUID",
			})
			return
		}

		_, err = db.Exec("INSERT INTO users(id,username, email, passhash, role) VALUES ($1,$2, $3, $4, $5)", newUser.ID, signUpUser.Username, signUpUser.Email, newUser.Passhash, signUpUser.Role)
		if err != nil {
			fmt.Println("Error inserting new user to DB", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error inserting new user to DB",
			})
			return
		}
		fmt.Println("New User Registration Successful")
		json.NewEncoder(w).Encode(Response{
			Code:    200,
			Message: "New User Registration Successful",
		})
		return

	case "sign-in":
		var loginDetails LogInUserDetails
		var dbUser DBUser

		err := json.NewDecoder(r.Body).Decode(&loginDetails)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		isEmailValid := ValidateEmail(loginDetails.Email)
		if !isEmailValid {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Email not in the correct format",
			})
			return
		}

		row := db.QueryRow("SELECT * FROM users where email= $1", loginDetails.Email)
		err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Passhash, &dbUser.Role)

		if err != nil && err != sql.ErrNoRows {
			fmt.Println("Error querying DB for users", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error querying DB for users",
			})
			return
		}

		//checks if email is registered
		if dbUser.Email == "" {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Email not registered. Please sign-up",
			})
			return
		}

		check := CheckPasswordHash(loginDetails.Password, dbUser.Passhash)

		if !check {
			fmt.Println("Username or Password is Incorrect")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Username or Password is Incorrect",
			})
			return
		}

		validToken, err := GenerateJWT(dbUser.Email, dbUser.Role)
		if err != nil {
			fmt.Println("Error generating JWT token", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error generating JWT token",
			})
			return
		}

		var token Token
		token.Email = dbUser.Email
		token.Role = dbUser.Role
		token.TokenString = validToken
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)

	case "mintNFT":
		JWTToken := r.Header.Get("Authorization")
		TokenArray := strings.Split(JWTToken, " ")
		if TokenArray[1] == "" {
			fmt.Println("No JWT Token found")
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "No JWT Token found",
			})
			return
		}
		claims, valid := ValidateToken(TokenArray[1])
		if !(valid) {
			fmt.Println("JWT Token error")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "JWT Token err",
			})
			return
		}

		if claims["role"] != "manufacturer" {
			fmt.Println("User role err")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "User role err. Only manufacturer can perform this action",
			})
			return
		}

		// err := MintToken(tokenURI)
		// if err != nil {
		// 	fmt.Println(err)
		// 	json.NewEncoder(w).Encode(Response{
		// 		Code:    500,
		// 		Message: "Internal server error",
		// 	})
		// 	return
		// }

		// json.NewEncoder(w).Encode(Response{
		// 	Code:    200,
		// 	Message: "Successfully minted NFT",
		// })

		// case "burnNFT":
		// JWTToken := r.Header.Get("Authorization")
		// TokenArray := strings.Split(JWTToken, " ")
		// if TokenArray[1] == "" {
		// 	fmt.Println("No JWT Token found")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    500,
		// 		Message: "No JWT Token found",
		// 	})
		// 	return
		// }
		// claims, valid := ValidateToken(TokenArray[1])
		// if !(valid) {
		// 	fmt.Println("JWT Token error")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    400,
		// 		Message: "JWT Token err",
		// 	})
		// 	return
		// }

		// if claims["role"] != "manufacturer" {
		// 	fmt.Println("User role err")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    400,
		// 		Message: "User role err. Only manufacturer can perform this action",
		// 	})
		// 	return
		// }

		// case "transferNFT":
		// 	JWTToken := r.Header.Get("Authorization")
		// 	TokenArray := strings.Split(JWTToken, " ")
		// 	if TokenArray[1] == "" {
		// 		fmt.Println("No JWT Token found")
		// 		json.NewEncoder(w).Encode(  Response{
		// 			Code:    500,
		// 			Message: "No JWT Token found",
		// 		})
		// 		return
		// 	}
		// 	_, valid := ValidateToken(TokenArray[1])
		// 	if !(valid) {
		// 		fmt.Println("JWT Token error")
		// 		json.NewEncoder(w).Encode(  Response{
		// 			Code:    400,
		// 			Message: "JWT Token err",
		// 		})
		// 		return
		// 	}
		//check to see if the caller of the function is the owner of the tokenuRI
	}
	return
}
