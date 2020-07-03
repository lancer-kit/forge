build:
	go get github.com/lancer-kit/forge

gen-default:
	./forge new -o="./project" -d="github.com" -n="holmanskih/project"

gen-simple:
	./forge new -o="./project" -n="holmanskih/project"

gen-with-gomod-api:
	./forge new -o="./project" -d="github.com" -n="holmanskih/project" --api

gen-with-gomod-db:
	./forge new -o="./project" -d="github.com" -n="holmanskih/project" --db

gen:
	./forge new

clean:
	rm -r ./project && rm ./forge

test:
	make build gen