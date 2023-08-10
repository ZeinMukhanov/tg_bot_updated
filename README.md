# tg_bot_updated

To run the code:

Fill out the database information in the code itself.

	password = "YOUR_PASSWORD"
	dbname   = "YOUR_DBNAME"
  botToken := "YOUR_BOTTOKEN"
 

Add the username and secret key for SIPUNI in sipuni.go

Terminal:

go build .

./tg_bot_updated

In PostgreSQL:
Create a table and add two columns:
1)email
2)name

Add the manager information to managers.csv:
email,name
