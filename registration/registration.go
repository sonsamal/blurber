package registration

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

type UserLedger interface {
	Add(u User, pwd string) error
	Remove(uname string, pwd string) error
	LogIn(uname string, pwd string) (error, string) // returns a token
	LogOut(uname string, pwd string) error
	CheckIn(uname string, token string) (error, string) // returns a token
}

type Token struct {
	creationDate time.Time
	token        string
}

type LocalUserLedger struct {
	userSet  map[int]User
	userID   map[string]int
	pwdMap   map[int]string
	tokenMap map[int]Token

	uidCounter int
}

func NewLocalLedger() *LocalUserLedger {
	return &LocalUserLedger{
		userSet:    make(map[int]User),
		userID:     make(map[string]int),
		pwdMap:     make(map[int]string),
		tokenMap:   make(map[int]Token),
		uidCounter: 0,
	}
}

// Not threadsafe
func (lul *LocalUserLedger) AddNewUser(name string, pwd string) error {
	log.Printf("REGISTRATION: Adding user %s with pwd %s", name, pwd)
	lul.userSet[lul.uidCounter] = User{
		Name: name,
		UID:  lul.uidCounter,
	}
	lul.userID[name] = lul.uidCounter
	lul.pwdMap[lul.uidCounter] = pwd

	lul.uidCounter++

	return nil
}

// Not threadsafe
func (lul *LocalUserLedger) LogIn(uname string, pwd string) (error, string) {
	log.Printf("LOGIN: Attempting user %s with pwd %s", uname, pwd)

	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), ""
	}

	if lul.pwdMap[id] != pwd {
		return errors.New("Bad password for " + uname), ""
	}

	return nil, lul.allocateNewToken(id)
}

func (lul *LocalUserLedger) allocateNewToken(id int) string {
	bitString := make([]byte, 256)
	_, err := rand.Read(bitString)
	if err != nil {
		panic(err)
	}

	token := hex.EncodeToString(bitString)

	lul.tokenMap[id] = Token{token: token, creationDate: time.Now()}

	return token
}

func (lul *LocalUserLedger) CheckIn(uname string, token string) (error, string) {
	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), ""
	}

	if lul.tokenMap[id].token != token {
		return errors.New("Bad token for " + uname), ""
	}

	if time.Since(lul.tokenMap[id].creationDate) > 30*time.Second {
		return errors.New("Session expried for " + uname), ""
	}

	return nil, lul.allocateNewToken(id)
}

func (lul *LocalUserLedger) GetUsrID(uname string) (error, int) {
	log.Printf("LOGIN: Getting user %s", uname)

	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), -1
	}
	return nil, id
}