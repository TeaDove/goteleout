.PHONY: ckeck install upload

check:
	pre-commit run -a
	go test -v ./...

install:
	go install

TAG := "${VERSION}.dev"
upload: # ONLY FOR DEV PURPOSES
	git add . || true
	git commit -m 'auto: uploading updates for tagging' || true
	git push
	git tag -d ${TAG} || true
	git tag -a ${TAG} -m "auto: development release"
	git push --delete origin ${TAG} || true
	git push origin ${TAG}
