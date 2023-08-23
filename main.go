package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type Profile struct {
	UUID       string
	Name       string
	Orders     []*Order
	Expiration time.Time
	sync.RWMutex
}

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

var TTL = time.Duration(time.Second * 2)

var profiles = make(map[string]Profile)

func main() {
	p := New()

	//fmt.Println(p)
	p2 := get(p.UUID)
	p.Insert("test")
	p2.Insert("test2")
	//p.Insert("test2")
	//p.Insert("test3")
	//fmt.Println(p)
	fmt.Println(p)
	fmt.Println(p2)
	fmt.Println(profiles)
	//for _, order := range p.Orders {
	//	fmt.Println(order.UUID, order.Value)
	//	p.Set(order.UUID, "TEST")
	//}
	//for i, order2 := range p.Orders {
	//	fmt.Println(i, order2)
	//}

}

func New() *Profile {

	var orders []*Order

	profile := Profile{
		Orders:     orders,
		Expiration: time.Now().Add(TTL),
		UUID:       uuid.NewString(),
	}

	profiles[profile.UUID] = profile

	return &profile
}

func (c *Profile) IsExpiredTTL() bool {
	if time.Now().After(c.Expiration) {
		*c = Profile{}
		return true
	}
	return false
}

func (c *Profile) Set(orderUUID string, value interface{}) *bool {

	if c.IsExpiredTTL() {
		return nil
	}
	c.Lock()
	defer c.Unlock()

	c.Expiration = time.Now().Add(TTL)

	idx := slices.IndexFunc(c.Orders, func(c *Order) bool { return c.UUID == orderUUID })
	if idx >= 0 {
		c.Orders[idx] = &Order{
			UUID:      orderUUID,
			Value:     value,
			CreatedAt: c.Orders[idx].CreatedAt,
			UpdatedAt: time.Now(),
		}
	} else {
		c.Orders = append(c.Orders, &Order{
			UUID:      orderUUID,
			Value:     value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	bool := true
	return &bool
}

func (c *Profile) Insert(value interface{}) *string {

	if c.IsExpiredTTL() {
		return nil
	}
	c.Lock()
	defer c.Unlock()
	id := uuid.NewString()
	c.Orders = append(c.Orders, &Order{
		UUID:      id,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	profiles[c.UUID] = *c
	return &id
}

func (c *Profile) Delete(orderUUID string) *bool {

	if c.IsExpiredTTL() {
		return nil
	}
	c.Lock()
	defer c.Unlock()

	idx := slices.IndexFunc(c.Orders, func(c *Order) bool { return c.UUID == orderUUID })

	if idx >= 0 {
		c.Orders = append(c.Orders[:idx], c.Orders[idx+1:]...)
		a := true
		return &a
	} else {
		b := false
		return &b
	}
}

func get(UUID string) *Profile {
	profile := profiles[UUID]
	return &profile
}
