PROTO_PATH = ./protos
GEN_PATH = ./gen

proto:
	@mkdir -p ${GEN_PATH}
	protoc -I ${PROTO_PATH} --go_out=plugins=grpc:${GEN_PATH} ${PROTO_PATH}/*.proto


%: ${PROTO_PATH}/%.proto
	@mkdir -p ${GEN_PATH}
	protoc -I ${PROTO_PATH} --go_out=plugins=grpc:${GEN_PATH} $^


clean:
	rm -f ${GEN_PATH}/*.pb.go
