run_test:
	go test -v -cover --coverprofile=demo.coverage -coverpkg=./... -count 1 github.com/vietky/golang-tmpl/demo/tests

check_coverage:
	go tool cover -html=demo.coverage

mock_gen:
	mockgen -source=demo/bettercodes/log.go  -destination=demo/mocks/log.go -package=mocks
