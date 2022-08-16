package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"
)

func TestBillForFunc(t *testing.T) {

	var expected float64

	// OUR EXPECTED NUMBER
	expected = float64(1.0938)	

	// OUR STARTING YEAR-MONTH
	str := "2022-08"

	var activeSubscription = Subscription{
		Id:             		1,
		CustomerId:          	1,
		MonthlyPriceInDollars:   5,
	}	

	// Open our jsonFile
	jsonFile, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	// fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	// data := User{}
	var data []User	

	err = json.Unmarshal([]byte(byteValue), &data)

	if err != nil { 
		fmt.Println(err)
		log.Fatal(err)
	}			

	actual := BillFor(str, &activeSubscription, &data)
	result := big.NewFloat(actual).Cmp(big.NewFloat(expected))

	if result != 0 {
		t.Errorf("expected %f, got %f", expected, actual)
	}

}

