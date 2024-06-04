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
	_, err = file.WriteString(`
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
	<script nonce="8aa92897-6d2f-481c-be30-dda0b4a9b403">try{(function(w,d){!function(j,k,l,m){j[l]=j[l]||{};j[l].executed=[];j.zaraz={deferred:[],listeners:[]};j.zaraz._v="5671";j.zaraz.q=[];j.zaraz._f=function(n){return async function(){var o=Array.prototype.slice.call(arguments);j.zaraz.q.push({m:n,a:o})}};for(const p of["track","set","debug"])j.zaraz[p]=j.zaraz._f(p);j.zaraz.init=()=>{var q=k.getElementsByTagName(m)[0],r=k.createElement(m),s=k.getElementsByTagName("title")[0];s&&(j[l].t=k.getElementsByTagName("title")[0].text);j[l].x=Math.random();j[l].w=j.screen.width;j[l].h=j.screen.height;j[l].j=j.innerHeight;j[l].e=j.innerWidth;j[l].l=j.location.href;j[l].r=k.referrer;j[l].k=j.screen.colorDepth;j[l].n=k.characterSet;j[l].o=(new Date).getTimezoneOffset();if(j.dataLayer)for(const w of Object.entries(Object.entries(dataLayer).reduce(((x,y)=>({...x[1],...y[1]})),{})))zaraz.set(w[0],w[1],{scope:"page"});j[l].q=[];for(;j.zaraz.q.length;){const z=j.zaraz.q.shift();j[l].q.push(z)}r.defer=!0;for(const A of[localStorage,sessionStorage])Object.keys(A||{}).filter((C=>C.startsWith("_zaraz_"))).forEach((B=>{try{j[l]["z_"+B.slice(7)]=JSON.parse(A.getItem(B))}catch{j[l]["z_"+B.slice(7)]=A.getItem(B)}}));r.referrerPolicy="origin";r.src="/cdn-cgi/zaraz/s.js?z="+btoa(encodeURIComponent(JSON.stringify(j[l])));q.parentNode.insertBefore(r,q)};["complete","interactive"].includes(k.readyState)?zaraz.init():j.addEventListener("DOMContentLoaded",zaraz.init)}(w,d,"zarazData","script");})(window,document)}catch(e){throw fetch("/cdn-cgi/zaraz/t"),e;};</script></head>
	<body class=" dark-mode hold-transition sidebar-mini">
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
	<img src="https://media.giphy.com/media/mAgG12Pk85e1mc31HJ/giphy.gif" alt="AdminLTE Logo" class="brand-image img-circle elevation-3" style="opacity: .8">
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
	<div class="row">
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
		<th>PhoneNmber</th>
		<th>Message</th>
		<th>Date</th>
	
	</tr>
	</thead>
	<tbody>`)
	if err != nil {
		return err
	}

	// Write user data to HTML table rows
	for i, user := range users {
		timestamp := user.Timestamp.Format("2006-01-02 15:04:05")
		// Tulis tag img ke file HTML
		var profilePhotoHTML string
		if user.ProfilePhotoURL != "" {
			profilePhotoHTML = "<a href='#' class='text-white'><img src='" + user.ProfilePhotoURL + "' alt='Profile Photo' width='50px' class='rounded-circle img-fluid'>  <span class='ps-2'>" + user.Username + "</span></a>"
		} else {
			profilePhotoHTML = "No photo"
		}
		_, err = file.WriteString("<tr> <td>" + strconv.Itoa(i+1) + "</td>  <td>" + profilePhotoHTML + "</td> <td>" + strconv.FormatInt(user.ID, 10) + "</td><td>" + user.FirstName + "</td><td>" + user.LastName + "</td><td>" + user.PhoneNumber + "</td><td>" + user.Message + "</td><td>" + timestamp + "</td></tr>")
		if err != nil {
			return err
		}
	}

	// Write HTML footer
	_, err = file.WriteString(`</tbody>
	</table>
	</div>
	</div>
	</div>
	</div>
	</div>
	</section>
	</div>
	<footer class="main-footer">
	<div class="float-right d-none d-sm-block">
	<b>Version</b> 0.0.2
	</div>
	<strong>Copyright &copy; 2024 <a href="https://aigoretech.rf.gd">BookfinderBot</a>.</strong> All rights reserved.
	</footer>
	<aside class="control-sidebar control-sidebar-dark">
	</aside>
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
	</body>
	</html>
	`)
	if err != nil {
		return err
	}

	return nil
}
