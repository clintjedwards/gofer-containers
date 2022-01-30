## build: build docker containers
build: check-semver-included
	docker build -f triggers/cron/Dockerfile -t ghcr.io/clintjedwards/gofer-containers/trigger_cron:${semver} .
	docker tag ghcr.io/clintjedwards/gofer-containers/trigger_cron:${semver} ghcr.io/clintjedwards/gofer-containers/trigger_cron:latest
	docker build -f triggers/interval/Dockerfile -t ghcr.io/clintjedwards/gofer-containers/trigger_interval:${semver} .
	docker tag ghcr.io/clintjedwards/gofer-containers/trigger_interval:${semver} ghcr.io/clintjedwards/gofer-containers/trigger_interval:latest
	docker build -f triggers/github/Dockerfile -t ghcr.io/clintjedwards/gofer-containers/trigger_github:${semver} .
	docker tag ghcr.io/clintjedwards/gofer-containers/trigger_github:${semver} ghcr.io/clintjedwards/gofer-containers/trigger_github:latest

## push: push docker to github
push: check-semver-included
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_cron:${semver}
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_cron:latest
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_interval:${semver}
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_interval:latest
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_github:${semver}
	docker push ghcr.io/clintjedwards/gofer-containers/trigger_github:latest

## help: prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

check-semver-included:
ifndef semver
	$(error semver is undefined; ex. semver=0.0.1)
endif

