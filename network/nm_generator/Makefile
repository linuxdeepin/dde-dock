all: gen-back-end-go

gen-back-end-go: gen-back-end-go-nm-setting-bean gen-back-end-go-nm-enums gen-back-end-go-nm-setting-const

gen-back-end-go-nm-enums:
	@./gen_nm_enums.sh

gen-back-end-go-nm-setting-const:
	@./gen_nm_setting_const.sh

gen-back-end-go-nm-setting-bean:
	go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --gen-back-end-go --back-end-go-dir ".."

gen-front-end-qt:
	@go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --gen-front-end-qt --front-end-qt-dir "/workspace/code/deepin/dde/dde-control-center/src/modules/network/edit_autogen/"

gen-front-end-qml:
	@go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --write --gen-front-end-qml --front-end-qml-dir "/workspace/code/deepin/dde/dde-control-center/src/modules/network/edit_autogen/"
