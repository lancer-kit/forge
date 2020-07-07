build:
	go get github.com/lancer-kit/forge

scaffold-in-gopath:
	./forge new --name github.com/out

scaffold-with-gomods:
	./forge new --name github.com/out --gomods out

scaffold-with-gopath-git:
	./forge new --name github.com/out --repo https://github.com/scaffolder/test.git

scaffold-with-gomods-git:
	scaffold-with-gomods --repo https://github.com/scaffolder/test.git
