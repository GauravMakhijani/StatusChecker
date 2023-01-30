# Website Monitor API

A REST API for monitoring the status of a list of websites.

## Getting Started

1. Clone the repository
2. install dependencies "go mod tidy"
3. setup database   
    
    $ createdb website_monitor

4.migrate
    migrate -database postgres://user:password@localhost:5432/dbname -path . up


