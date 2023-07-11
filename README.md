sudo snap install sqlc

sqlc init #create sqlc.yaml file

sqlc generate


for test running

go test -timeout 30s ./db/sqlc/ -run TestMain

go get github.com/stretchr/testify



use dbdiagrams to create diagrams and download the sql file
use go migrate tool for migration
use sqlc to generate go code from sql code
use testify to test the code
