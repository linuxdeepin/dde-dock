GO15VENDOREXPERIMENT=0

.PHONY: nm_consts_gen.yml
nm_consts_gen.yml:
	./gen_nm_consts.py --gir /usr/share/gir-1.0/NM-1.0.gir --output ./nm_consts_gen.yml

gen-nm-code: nm_consts_gen.yml
	mkdir -p ../nm
	go run $(shell find -type f -iname "*.go" -not -iname "*test.go") --output-consts ../nm/nm_consts_gen.go --output-beans ../nm_setting_beans_gen.go
	gofmt -w ../nm/nm_consts_gen.go
	gofmt -w ../nm_setting_beans_gen.go
