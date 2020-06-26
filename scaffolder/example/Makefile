build:
	go get github.com/lancer-kit/forge

gen-default:
	./forge gen -o="./project" -d="github.com" -n="holmanskih/project"

gen-simple:
	./forge gen -o="./project" -n="holmanskih/project"

gen-with-gomod-api:
	./forge gen -o="./project" -d="github.com" -n="holmanskih/project" --api

gen-with-gomod-db:
	./forge gen -o="./project" -d="github.com" -n="holmanskih/project" --db

gen:
	./forge gen

clean:
	rm -r ./project && rm ./forge

test:
	make build gen