all:
	go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --back-end --back-end-dir ".." --front-end --front-end-dir "$$HOME/workspace/code/deepin/dde-control-center/src/modules/network/edit_autogen/"

back-end:
	go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --back-end --back-end-dir ".."

front-end:
	go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --front-end --front-end-dir "$$HOME/workspace/code/deepin/dde-control-center/src/modules/network/edit_autogen/"
