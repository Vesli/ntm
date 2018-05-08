package main

/*
   Remember to structure via data.
   ex:
   	Is there an email to send at subscription (email package)
	Is there any permission on event creation, subscription? (permissions package)

	USE Chi Render for responses and Content-Type JSON
*/

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	s := newService()
	defer s.close()
	s.run()
}
