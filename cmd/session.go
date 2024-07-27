package main

import (
	"github.com/alexedwards/scs/v2"
)

type CustomSessionManager struct {
	*scs.SessionManager
}

func NewCustomSessionManager() *CustomSessionManager {
	manager := scs.New()
	myManager := &CustomSessionManager{manager}
	return myManager
}
