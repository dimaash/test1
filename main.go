package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

type Subscription struct {
  Id int
  CustomerId int
  MonthlyPriceInDollars int
}

type User struct {
  Id int
  Name string
  ActivatedOn time.Time
  DeactivatedOn time.Time
  CustomerId int
}
func FirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
func LastDayOfMonth(t time.Time) time.Time {
	return FirstDayOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}
func NextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}
func IsActiveUser(user User) bool {

	// if !user.DeactivatedOn.IsZero() {
	// 	fmt.Printf("[WARNING] User Id: %d => DeactivatedOn on [%v] is ZERO Date\n", user.Id, user.DeactivatedOn.Format("02-Jan-2006"))
	// }

	if !user.DeactivatedOn.IsZero() && user.DeactivatedOn.Before(user.ActivatedOn) {
		fmt.Printf("[WARNING] User Id: %d => DeactivatedOn on [%v] is Before ActivatedOn on [%v]\n", user.Id, user.DeactivatedOn.Format("02-Jan-2006"), user.ActivatedOn.Format("02-Jan-2006"))
		return false
	}
	if user.ActivatedOn.IsZero() {
		fmt.Printf("[WARNING] User Id: %d => ActivatedOn is ZERO Date\n", user.Id)
		return false
	}
	if !user.DeactivatedOn.IsZero() && user.DeactivatedOn.Before(time.Now()) {
		fmt.Printf("[WARNING] User Id: %d => DeactivatedOn on [%v] is Before NOW on [%v]\n", user.Id, user.DeactivatedOn.Format("02-Jan-2006"), time.Now().Format("02-Jan-2006"))
		return false
	}	

    return user.ActivatedOn.Before(time.Now())
}

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}


func BillFor(yearMonth string, activeSubscription *Subscription, users *[]User) float64 {

	if activeSubscription == nil {
		return 0
	}
	if len(*users) == 0 {
		return 0
	}

	format := "2006-01"

	dateTime, err := time.Parse(format, yearMonth)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println(dateTime)		

	t1 := LastDayOfMonth(dateTime)
	t2 := FirstDayOfMonth(dateTime)
	t1minus1 := t2.AddDate(0, 0, -1)

	difference := (t1.Sub(t2).Hours() / 24) + 1

	fmt.Printf("%d Days in this month of %v\n", int64(difference), dateTime.Format("Jan"))

	tnextDay := NextDay(t1minus1)

	fmt.Printf("Monthly Price in Dollars $: %f\n", float64(activeSubscription.MonthlyPriceInDollars))

	var monthlyRate = float64(activeSubscription.MonthlyPriceInDollars) / difference

	fmt.Printf("Monthly Rate in Dollars $: %f\n", float64(monthlyRate))

	var subtoal float64
	subtoal = float64(0)

	// LETS LOOP THROUGH ALL THE DAYS IN THE GIVEN MONTH
	for i := 0; i < int(difference); i++ {

		var tempUsers []User

		// LETS LOOP THROUGH OUR SLICE OF USERS THAT WE HAVE READ FROM users.json
		for _, user := range *users {

			// IF THE USER IS ACTIVE AND its ActiveDate is equals to THE DAY OF THE MONTH IN QUESTION
			// LETS ADD IT TO THE TEMPORARY STORAGE
			if (IsActiveUser(user) && user.ActivatedOn.Truncate(24*time.Hour).Equal(tnextDay.Truncate(24*time.Hour)) ) {
				// fmt.Printf("User Id: %d IS ACTIVE\n", user.Id)
				tempUsers = append(tempUsers, user)
			} 
		}

		// Multiply the number of active users for the day by the daily rate to calculate the total for the day
		subtoal += float64(len(tempUsers)) * monthlyRate
		
		tnextDay = NextDay(tnextDay)
	}	

	subtoal = roundFloat(subtoal, 4) //  math.Floor(subtoal*100)/100

	return subtoal
}


func main() {

    pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(pwd)

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

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	// json.Unmarshal(byteValue, &data)
	err = json.Unmarshal([]byte(byteValue), &data)

	if err != nil { 
		fmt.Println(err)
		log.Fatal(err)
	}

	var subtoal float64
	subtoal = float64(0)

	subtoal = BillFor(str, &activeSubscription, &data)

	fmt.Printf("Roilling Subtoal = $%f\n", subtoal)

	fmt.Println("==========> DONE <===========")
}