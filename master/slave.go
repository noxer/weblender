package master

import (
	"time"
)

// =================================== STRUCT =================================

type slave struct {
	name     string
	lastSeen time.Time
}

// =================================== CONSTRUCTOR ============================

func NewSlave(name string) *slave {
	return &slave{
		name:     name,
		lastSeen: time.Now(),
	}
}

// =================================== GETTER =================================

func (s *slave) Name() string {
	return s.name
}

func (s *slave) LastSeen() time.Time {
	return s.lastSeen
}

// =================================== SETTER =================================

func (s *slave) Seen() {
	s.lastSeen = time.Now()
}
