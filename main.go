package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Booking struct {
	Name      string `json:"name"`
	CheckIn   string `json:"checkin"`
	CheckOut  string `json:"checkout"`
	RoomType  string `json:"roomtype"`
	Email     string `json:"email"`
	BookingID string `json:"booking_id"`
}

var bookings []Booking

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	loadBookings()

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/booking", BookingPage)
	http.HandleFunc("/book", HandleBooking)
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	err = tmpl.Execute(w, bookings)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func BookingPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/hotel.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	err = tmpl.Execute(w, bookings)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func loadBookings() {
	file, err := os.Open("bookings.json")
	if err != nil {
		log.Println("No existing bookings file found. Starting fresh.")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading bookings file: %v", err)
	}

	err = json.Unmarshal(data, &bookings)
	if err != nil {
		log.Fatalf("Error parsing bookings file: %v", err)
	}
}

func saveBookings() {
	data, err := json.MarshalIndent(bookings, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding bookings to JSON: %v", err)
	}

	err = ioutil.WriteFile("bookings.json", data, 0644)
	if err != nil {
		log.Fatalf("Error writing bookings to file: %v", err)
	}
}

func generateBookingID() string {
	rand.Seed(time.Now().UnixNano()) // Seed untuk memastikan angka acak berbeda setiap kali dijalankan
	return fmt.Sprintf("%07d", rand.Intn(100000))
}

func HandleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		booking := Booking{
			Name:      r.FormValue("name"),
			CheckIn:   r.FormValue("checkin"),
			CheckOut:  r.FormValue("checkout"),
			RoomType:  r.FormValue("roomtype"),
			Email:     r.FormValue("email"),
			BookingID: generateBookingID(),
		}
		bookings = append(bookings, booking)

		saveBookings()

		http.Redirect(w, r, "/booking", http.StatusSeeOther)
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
