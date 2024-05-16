TARGETBIN=ithings
.PHONY:	all ${TARGETBIN}.exe ${TARGETBIN} protoc

BUILD_ROOT=$(PWD)
all: ${TARGETBIN}.exe  ${TARGETBIN}

${TARGETBIN}:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.cn && \
	go build -ldflags "-w -s" -o $@ ithings.go
	@chmod 777 $@
	

protoc:
	@rm -rf ${BUILD_ROOT}/transport/isync/isync*.pb.go
	@protoc --go_out=${BUILD_ROOT}/transport/isync/ --go_opt=paths=source_relative \
	--go-grpc_out=${BUILD_ROOT}/transport/isync/ --go-grpc_opt=paths=source_relative \
	--proto_path=${BUILD_ROOT}/transport/isync/ \
	${BUILD_ROOT}/transport/isync/isync.proto

install:
	@mkdir -p out
	@chmod 777 ${TARGETBIN}.exe  ${TARGETBIN}
	@cp -a conf ${TARGETBIN}.exe  ${TARGETBIN}  out/
	sync;sync
	@echo "[Done]"

.PHONY: clean  install
clean:
	@rm -rf ${TARGETBIN}.exe  ${TARGETBIN} *.log *.db *.tar.gz
	@echo "[clean Done]"
