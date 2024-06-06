package datauser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// UserData represents data of a user
type UserData struct {
	ID              int64     `json:"id"`
	ProfilePhotoURL string    `json:"profile_photo_url"`
	Username        string    `json:"username"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PhoneNumber     string    `json:"phone_number"`
	Messages        []Message `json:"messages"`
}

// Message represents a single message in the conversation
type Message struct {
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

// LoadUserData loads user data from a JSON file
func LoadUserData(filename string) ([]UserData, error) {
	var users []UserData
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil // Return empty slice if file does not exist
		}
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// SaveUserData saves user data to a JSON file
func SaveUserData(filename string, users []UserData) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AddUserMessage adds a new message to the user's conversation
func AddUserMessage(users []UserData, userID int64, message Message) ([]UserData, error) {
	for i, user := range users {
		if user.ID == userID {
			users[i].Messages = append(users[i].Messages, message)
			return users, nil
		}
	}

	return nil, os.ErrNotExist
}

// SaveUserDataToHTML saves user data to an HTML file
func SaveUserDataToHTML(users []UserData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write HTML header
	header := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>BookFinderBot | DataTables</title>
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,400i,700&display=fallback">
<link rel="stylesheet" href="https://adminlte.io/themes/v3/plugins/fontawesome-free/css/all.min.css">
<link rel="stylesheet" href="https://adminlte.io/themes/v3/plugins/datatables-bs4/css/dataTables.bootstrap4.min.css">
<link rel="stylesheet" href="https://adminlte.io/themes/v3/plugins/datatables-responsive/css/responsive.bootstrap4.min.css">
<link rel="stylesheet" href="https://adminlte.io/themes/v3/plugins/datatables-buttons/css/buttons.bootstrap4.min.css">
<link rel="stylesheet" href="https://adminlte.io/themes/v3/dist/css/adminlte.min.css?v=3.2.0">
<script nonce="8aa92897-6d2f-481c-be30-dda0b4a9b403">try{(function(w,d){!function(j,k,l,m){j[l]=j[l]||{};j[l].executed=[];j.zaraz={deferred:[],listeners:[]};j.zaraz._v="5671";j.zaraz.q=[];j.zaraz._f=function(n){return async function(){var o=Array.prototype.slice.call(arguments);j.zaraz.q.push({m:n,a:o})}};for(const p of["track","set","debug"])j.zaraz[p]=j.zaraz._f(p);j.zaraz.init=()=>{var q=k.getElementsByTagName(m)[0],r=k.createElement(m),s=k.getElementsByTagName("title")[0];s&&(j[l].t=k.getElementsByTagName("title")[0].text);j[l].x=Math.random();j[l].w=j.screen.width;j[l].h=j.screen.height;j[l].j=j.innerHeight;j[l].e=j.innerWidth;j[l].l=j.location.href;j[l].r=k.referrer;j[l].k=j.screen.colorDepth;j[l].n=k.characterSet;j[l].o=(new Date).getTimezoneOffset();if(j.dataLayer)for(const w of Object.entries(Object.entries(dataLayer).reduce(((x,y)=>({...x[1],...y[1]})),{})))zaraz.set(w[0],w[1],{scope:"page"});j[l].q=[];for(;j.zaraz.q.length;){const z=j.zaraz.q.shift();j[l].q.push(z)}r.defer=!0;for(const A of[localStorage,sessionStorage])Object.keys(A||{}).filter((C=>C.startsWith("_zaraz_"))).forEach((B=>{try{j[l]["z_"+B.slice(7)]=JSON.parse(A.getItem(B))}catch{j[l]["z_"+B.slice(7)]=A.getItem(B)}}));r.referrerPolicy="origin";r.src="/cdn-cgi/zaraz/s.js?z="+btoa(encodeURIComponent(JSON.stringify(j[l])));q.parentNode.insertBefore(r,q)};["complete","interactive"].includes(k.readyState)?zaraz.init():j.addEventListener("DOMContentLoaded",zaraz.init)}(w,d,"zarazData","script");})(window,document)}catch(e){throw fetch("/cdn-cgi/zaraz/t"),e;};</script>
</head>
<body class="dark-mode hold-transition sidebar-mini">
<div class="wrapper">
<nav class="main-header navbar navbar-expand navbar-white navbar-light">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a>
</li>
</ul>
<ul class="navbar-nav ml-auto">
<li class="nav-item">
<a class="nav-link" data-widget="navbar-search" href="#" role="button">
<i class="fas fa-search"></i>
</a>
<div class="navbar-search-block">
<form class="form-inline">
<div class="input-group input-group-sm">
<input class="form-control form-control-navbar" type="search" placeholder="Search" aria-label="Search">
<div class="input-group-append">
<button class="btn btn-navbar" type="submit">
<i class="fas fa-search"></i>
</button>
<button class="btn btn-navbar" type="button" data-widget="navbar-search">
<i class="fas fa-times"></i>
</button>
</div>
</div>
</form>
</div>
</li>
<li class="nav-item">
<a class="nav-link" data-widget="fullscreen" href="#" role="button">
<i class="fas fa-expand-arrows-alt"></i>
</a>
</li>
<li class="nav-item">
<a class="nav-link" data-widget="control-sidebar" data-slide="true" href="#" role="button">
<i class="fas fa-th-large"></i>
</a>
</li>
</ul>
</nav>
<aside class="main-sidebar sidebar-dark-primary elevation-4">
<a href="https://aigoretech.rf.gd" class="brand-link">
<img src="https://media.giphy.com/media/mAgG12Pk85e1mc31HJ/giphy.gif" alt="BookFinderBot Logo" class="brand-image img-circle elevation-3" style="opacity: .8">
<span class="brand-text font-weight-light">BookFinderBot</span>
</a>
<div class="sidebar">
<nav class="mt-2">
<ul class="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">
<li class="nav-item">
<a href="#" class="nav-link">
<i class="nav-icon fas fa-tachometer-alt"></i>
<p>Dashboard</p>
</a>
</li>
<li class="nav-item">
<a href="#" class="nav-link">
<i class="nav-icon fas fa-table"></i>
<p>Tables</p>
</a>
</li>
</ul>
</nav>
</div>
</aside>
<div class="content-wrapper">
<section class="content-header">
<div class="container-fluid">
<div class="row mb-2">
<div class="col-sm-6">
<h1>DataTables</h1>
</div>
<div class="col-sm-6">
<ol class="breadcrumb float-sm-right">
<li class="breadcrumb-item"><a href="#">Home</a></li>
<li class="breadcrumb-item active">DataTables</li>
</ol>
</div>
</div>
</div>
</section>
<section class="content">
<div class="container-fluid">
<div class="row">`
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	table := `
<div class="col-12">
<div class="card">
<div class="card-header">
<h3 class="card-title">DataTable Users</h3>
</div>
<div class="card-body">
<table id="example1" class="table table-bordered table-striped">
<thead>
<tr>
<th>#</th>
<th>Profile</th>
<th>ID</th>
<th>FirstName</th>
<th>LastName</th>
<th>PhoneNumber</th>
<th>Messages</th>
<th>Date</th>
</tr>
</thead>
<tbody>`
	_, err = file.WriteString(table)
	if err != nil {
		return err
	}

	// Write user data to table
	for i, user := range users {
		var lastMessage string
		if len(user.Messages) > 0 {
			lastMessage = user.Messages[len(user.Messages)-1].Content
		} else {
			lastMessage = "No messages"
		}

		timestamp := user.Messages[len(user.Messages)-1].Timestamp.Format("2006-01-02 15:04:05")
		var profilePhotoHTML string
		if user.ProfilePhotoURL != "" {
			profilePhotoHTML = "<a href='#" + user.Username + "' class='text-white nav-link' data-toggle='tab'><img src='" + user.ProfilePhotoURL + "' alt='Profile Photo' width='50px' class='rounded-circle img-fluid'>" + user.Username + "</a>"
		} else {
			profilePhotoHTML = "No photo"
		}
		_, err = file.WriteString("<tr><td>" + strconv.Itoa(i+1) + "</td><td>" + profilePhotoHTML + "</td><td>" + strconv.FormatInt(user.ID, 10) + "</td><td>" + user.FirstName + "</td><td>" + user.LastName + "</td><td>" + user.PhoneNumber + "</td><td>" + lastMessage + "</td><td>" + timestamp + "</td></tr>")
		if err != nil {
			return err
		}
	}

	closeTable := `
</tbody>
</table>
</div>
</div>
</div>`
	_, err = file.WriteString(closeTable)
	if err != nil {
		return err
	}

	// Direct Chat section
	directChat := `
<div class="col-12">
<div class="card direct-chat direct-chat-primary">
<div class="card-header ui-sortable-handle">
<h3 class="card-title">Direct Chat</h3>
<div class="card-tools">
<button type="button" class="btn btn-tool" data-card-widget="collapse">
<i class="fas fa-minus"></i>
</button>
<button type="button" class="btn btn-tool" title="Contacts" data-widget="chat-pane-toggle">
<i class="fas fa-comments"></i>
<span class="badge badge-primary navbar-badge">3</span>
</button>
</div>
</div>
<div class="card-body">
<div class="tab-content">`
	_, err = file.WriteString(directChat)
	if err != nil {
		return err
	}

	// Write user messages and bot responses to direct chat
	for _, user := range users {
		_, err = file.WriteString("<div class='tab-pane' id='" + user.Username + "'><div class='direct-chat-messages'>")
		if err != nil {
			return err
		}

		// Write each message in the conversation
		for _, message := range user.Messages {
			var senderName string
			var floatClass string
			var senderProfilePhotoURL string
			if message.Sender == "user" {
				senderName = user.Username
				floatClass = "float-left"
				senderProfilePhotoURL = user.ProfilePhotoURL
			} else if message.Sender == "bot" {
				senderName = "BookFinderBot"
				floatClass = "float-right"
				// Ganti URL gambar default dengan URL gambar bot yang Anda berikan
				senderProfilePhotoURL = "https://media.giphy.com/media/mAgG12Pk85e1mc31HJ/giphy.gif"

			}

			if message.Sender == "bot" {
				floatClass = "float-right"
			}

			_, err = file.WriteString("<div class='direct-chat-msg'><div class='direct-chat-infos clearfix'><span class='direct-chat-name " + floatClass + "'>" + senderName + "</span><span class='direct-chat-timestamp " + floatClass + "'>" + message.Timestamp.Format("2006-01-02 15:04:05") + "</span></div><img class='direct-chat-img' src='" + senderProfilePhotoURL + "' alt='message " + message.Sender + " image'><div class='direct-chat-text'>" + message.Content + "</div></div>")
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString("</div></div>")
		if err != nil {
			return err
		}
	}
	// Close Direct Chat
	closeDC := `
</div>
</div>
</div>
</div>`
	_, err = file.WriteString(closeDC)
	if err != nil {
		return err
	}

	// Footer
	footer := `
</div>
</div>
</section>
</div>
<footer class="main-footer">
<div class="float-right d-none d-sm-block">
<b>Version</b> 3.2.0
</div>
<strong>Copyright &copy; 2014-2021 <a href="https://adminlte.io">AdminLTE.io</a>.</strong> All rights reserved.
</footer>
<aside class="control-sidebar control-sidebar-dark"></aside>
</div>
<script src="https://adminlte.io/themes/v3/plugins/jquery/jquery.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/bootstrap/js/bootstrap.bundle.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables/jquery.dataTables.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-bs4/js/dataTables.bootstrap4.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-responsive/js/dataTables.responsive.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-responsive/js/responsive.bootstrap4.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-buttons/js/dataTables.buttons.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-buttons/js/buttons.bootstrap4.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/jszip/jszip.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/pdfmake/pdfmake.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/pdfmake/vfs_fonts.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-buttons/js/buttons.html5.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-buttons/js/buttons.print.min.js"></script>
<script src="https://adminlte.io/themes/v3/plugins/datatables-buttons/js/buttons.colVis.min.js"></script>
<script src="https://adminlte.io/themes/v3/dist/js/adminlte.min.js?v=3.2.0"></script>
<script>
$(function () {
$("#example1").DataTable({
"responsive": true, "lengthChange": false, "autoWidth": false,
"buttons": ["copy", "csv", "excel", "pdf", "print", "colvis"]
}).buttons().container().appendTo('#example1_wrapper .col-md-6:eq(0)');
$('#example2').DataTable({
"paging": true,
"lengthChange": false,
"searching": false,
"ordering": true,
"info": true,
"autoWidth": false,
"responsive": true,
});
});
</script>
<script>
  $(document).ready(function(){
    // Toggle tab content on click
    $('.nav-link').on('click', function(){
      // Hapus kelas 'active' dari tab sebelumnya
      $('.nav-link').removeClass('active');
      // Sisipkan kelas 'active' pada tab yang sedang diklik
      $(this).addClass('active');
      // Ambil target konten tab yang akan ditampilkan
      var target = $(this).attr('href');
      // Sembunyikan semua tab content
      $('.tab-pane').removeClass('active');
      $('.tab-pane').hide();
      // Tampilkan tab content yang sesuai dengan tab yang sedang aktif
      $(target).addClass('active');
      $(target).show();
    });
  });
</script>
</body>
</html>`
	_, err = file.WriteString(footer)
	if err != nil {
		return err
	}

	return nil
}
