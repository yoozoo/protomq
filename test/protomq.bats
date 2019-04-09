#!/usr/bin/env bats

@test "test.proto go output" {
	../protomq gen --lang=goproducer result/go/test/producer proto/test.proto
	../protomq gen --lang=goconsumer result/go/test/consumer proto/test.proto

	diff -I "^//.*$" -x "*.pb.go" -r result/go/ expected/go/
}

@test "test.proto php output" {
	../protomq gen --lang=phpconsumer result/php/test/consumer proto/test.proto
	../protomq gen --lang=phpproducer result/php/test/producer proto/test.proto
	diff -I "^//.*$" -r -n result/php/ expected/php/
}
