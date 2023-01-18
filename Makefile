MOCKPKG = mock

OUT = test_out

test:
	mockgen -source=internal/form/repo.go -destination=internal/form/mock/pg_repo_mock.go -package=$(MOCKPKG)
	mockgen -source=internal/question/repo.go -destination=internal/question/mock/pg_repo_mock.go -package=$(MOCKPKG)
	mockgen -source=internal/answer/repo.go -destination=internal/answer/mock/pg_repo_mock.go -package=$(MOCKPKG)
	mockgen -source=internal/poolanswer/repo.go -destination=internal/poolanswer/mock/pg_repo_mock.go -package=$(MOCKPKG)
	mockgen -source=internal/auth/repo.go -destination=internal/auth/mock/pg_repo_mock.go -package=$(MOCKPKG)
	mkdir -p $(OUT)/
	go test ./internal/form/usecase ./internal/form/repo \
	./internal/question/usecase ./internal/question/repo \
	./internal/answer/usecase ./internal/answer/repo \
	./internal/poolanswer/usecase ./internal/poolanswer/repo \
	./internal/auth/usecase ./internal/auth/repo \
	-v -cover -coverprofile=$(OUT)/coverage.out >> $(OUT)/report.txt
	go tool cover -html=$(OUT)/coverage.out -o $(OUT)/index.html

clean:
	rm -rf internal/form/mock
	rm -rf internal/question/mock
	rm -rf internal/answer/mock
	rm -rf internal/auth/mock
	rm -rf internal/poolanswer/mock
	rm -rf $(OUT)
