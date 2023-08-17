# Setting Up and Testing the Random Quote and Image Application

This document provides instructions on how to set up and test the Random Quote and Image Application. The application is divided into two parts: 'terminal' and 'web'. Follow the steps below to get the application up and running.

## Prerequisites

Make sure you have the following prerequisites installed on your system:
- Go programming language (version 1.21)
- Git

## Clone the Repository

Open a terminal and clone the repository to your desired directory:
`git clone <https://github.com/ramyadmz/tucows>`
`cd <repository-directory>`

## Terminal Application

### Building and Running the Terminal Application

Navigate to the 'cmd/terminal' directory:
`cd cmd/terminal`

Build the terminal application:
`go build -o terminal-app`

Run the terminal application:
`./terminal-app`

### Testing the Terminal Application

To test the terminal application, you can use the following flags:
- '-category': Specify the quote category (optional)
- '-width': Specify the image width (default: 40)
- '-height': Specify the image height (default: 30)
- '-filters': Specify image filters as a comma-separated list (e.g., "grayscale,blur")

Example command with flags:
`./terminal-app -category 1 -width 80 -height 60 -filters grayscale,blur`

## Web Application

### Building and Running the Web Application

Navigate to the 'cmd/webapp' directory:
`cd cmd/webapp`

Build the web application:
`go build -o web-app`

Run the web application:
`./web-app`

You can use the following flags:
- '-port': Specify the localhost port of our web app (optional)

### Testing the Web Application

The web application can be accessed by opening a web browser and navigating to 'http://localhost:8080'.

You can add query parameters to the URL to customize the behavior:
- 'key': Specify the quote category (optional)
- 'width': Specify the image width (default: 600)
- 'height': Specify the image height (default: 400)
- 'filters': Specify image filters as a comma-separated list (e.g., "grayscale,blur")

Example URL with query parameters:
'http://localhost:8080?key=1&width=800&height=600&filters=grayscale,blur'