package datauser

import (
	"os"
	"strconv"
	"time"
)

// UserData represents data of a user
type UserData struct {
	ID              int64
	ProfilePhotoURL string
	Username        string
	FirstName       string
	LastName        string
	PhoneNumber     string
	Message         string
	Timestamp       time.Time
	// Add other fields as needed
}

// SaveUserDataToHTML saves user data to an HTML file
func SaveUserDataToHTML(users []UserData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write HTML header
	_, err = file.WriteString(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Table User</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
  </head>
  <body>
  <div class="table-responsive">
    <table class="table table-bordered border-primary" border="1">
     <thead class="table-dark">
      <tr>
       <th scope="col">No</th>
			 <th scope="col">Profile</th>
        <th scope="col">ID</th>
       	<th scope="col">FirstName</th>
	<th scope="col">LastName</th>
	<th scope="col">PhoneNumber</th>
        <th scope="col">Message</th>
	<th scope="col">Timestamp</th>
	

      </tr>
        </thead>`)
	if err != nil {
		return err
	}

	// Write user data to HTML table rows
	for i, user := range users {
		timestamp := user.Timestamp.Format("2006-01-02 15:04:05")
		// Tulis tag img ke file HTML
		var profilePhotoHTML string
		if user.ProfilePhotoURL != "" {
			profilePhotoHTML = "<img src='" + user.ProfilePhotoURL + "' alt='Profile Photo' width='50' height='50' class='rounded-circle img-fluid'>  <span class='ps-2'>" + user.Username + "</span>"
		} else {
			profilePhotoHTML = "No photo"
		}
		_, err = file.WriteString("<tr> <td>" + strconv.Itoa(i+1) + "</td>  <td>" + profilePhotoHTML + "</td> <td>" + strconv.FormatInt(user.ID, 10) + "</td><td>" + user.FirstName + "</td><td>" + user.LastName + "</td><td>" + user.PhoneNumber + "</td><td>" + user.Message + "</td><td>" + timestamp + "</td></tr>")
		if err != nil {
			return err
		}
	}

	// Write HTML footer
	_, err = file.WriteString(`    </table></div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
  </body>
</html>`)
	if err != nil {
		return err
	}

	return nil
}
