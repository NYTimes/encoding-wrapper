all: test

testdeps:
	go get -d -t ./...

checkfmt: testdeps
	[ -z "$$(gofmt -s -d . | tee /dev/stderr)" ]

deadcode:
	go get github.com/remyoudompheng/go-misc/deadcode
	go list ./... | sed -e "s;github.com/NYTimes/encoding-wrapper/;;" | xargs deadcode

lint: testdeps
	go get github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go'); do \
		export output="$$(golint $${file})"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

test: checkfmt lint vet deadcode
	go test ./...

vet: testdeps
	go vet ./...
