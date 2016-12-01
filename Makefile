SUBPACKAGES=. rfc3164 rfc5424

help:
	@echo "Available targets"
	@echo "- tests: run tests"
	@echo "- installdependencies: installs dependencies declared in dependencies.txt"
	@echo "- benchmarks: run benchmarks"

installdependencies:
	@cat dependencies.txt | xargs go get

tests: installdependencies
	@for pkg in $(SUBPACKAGES); do cd $$pkg && go test -v -i && go test -v; cd -;done

benchmarks:
	@for pkg in $(SUBPACKAGES); do cd $$pkg && go test -gocheck.b ; cd -;done
