#!/usr/bin/env bats

@test "test.proto go output" {
	../protomq gen --lang=go result/go/test proto/test.proto

	diff -I "^//.*$" -x "*.pb.go" -r result/go/ expected/go/
}

@test "test.proto php output" {
	skip
	../protomq gen --lang=php result/php/test proto/test.proto
	diff -I "^//.*$" -r -n result/php/ expected/php/
}
