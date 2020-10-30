CACHE_DIR       = .cache/coverage.out
GOCOV_HTML_FILE = test-coverage.html

test:
	@mkdir -p "$(CACHE_DIR)"
	@ echo "go test -cover"
	@ go test -v -covermode=count "-coverprofile=$(CACHE_DIR)/cov.out"
	@ go tool cover "-html=$(CACHE_DIR)/cov.out" -o "$(CACHE_DIR)/cov.html"
	@  sed 's/.cov0 { color: rgb(192, 0, 0)/.cov0 { color: rgb(255, 100, 80)/g' \
	   "$(CACHE_DIR)/cov.html" \
	 | sed 's/font-weight: bold/font-weight: normal/g' \
	 | sed 's/font-family:/tab-size:2;font-family: SFMono-Regular,Consolas,Liberation Mono,Menlo,/g'\
	 | sed 's/background: black;/background: rgba(20,20,20);/g' \
	 | python -c 'import re,sys;print(re.sub(r"\n {8}", "\n\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 | python -c 'import re,sys;print(re.sub(r"\n(\t+) {8}", "\n\\1\t", sys.stdin.read()))' \
	 > "$(GOCOV_HTML_FILE)"
	 @ echo "test coverage report written to $(GOCOV_HTML_FILE)"

fmt:
	gofmt -w -s -l .

doc:
	@echo "open http://localhost:6060/pkg/github.com/rsms/go-uuid/"
	@bash -c '[ "$$(uname)" == "Darwin" ] && \
	         (sleep 1 && open "http://localhost:6060/pkg/github.com/rsms/go-uuid/") &'
	godoc -http=localhost:6060

clean:
	rm -rvf "$(GOCOV_HTML_FILE)" "$(CACHE_DIR)"

.PHONY: test clean release dist fmt doc dev dev1
