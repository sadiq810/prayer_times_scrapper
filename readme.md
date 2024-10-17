## Fetch Countries, Cities, Masjid and Masjid Prayer Timings from MyMasjid.com

This project is developed in golang & using goroutines for saving & processing data fetched form the APIs.

### Installation

- Create db and import DB_DUMP.sql file.
- Set DB information in .env file.
- To Run the code just run the following command:

- ``go run main.go``

### To Create Build for different platforms

``go build -o out/put/directory``
- Check you env regarding your architecture.

`` go env ``

- for more details
``go tool dist list``

### To create an executable for different platform
- Open cmd with administrative privileges on windows and run the following
1) set GOARCH=amd64
2) set GOOS=linux
3) go build 