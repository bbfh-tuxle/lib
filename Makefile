test:
	@go test -v ./tests/...

coverage:
	@go test -cover ./tuxle/... ./tests/... \
	| grep "tests" \
	| sed 's/\t/    /g' \
	| awk '{gsub(/([0-9]+\.[0-9]+%)/, "\033[1;31m&\033[0m"); print}'

