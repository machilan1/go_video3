package main

import (
	"context"
	"encoding/json"
	"fmt"

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

func (sm *CustomSessionManager) PutStruct(c context.Context, key string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	sm.Put(c, key, bytes)
	return nil
}

func (sm *CustomSessionManager) PopStruct(c context.Context, key string, dest any) error {

	bytes := sm.Pop(c, key)
	switch bytes.(type) {
	case []byte:
		json.Unmarshal(bytes.([]byte), dest)
		return nil
	default:
		return fmt.Errorf("cannot pop struct from key %s", key)
	}

}
