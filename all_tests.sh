#!/usr/bin/env bash
# Run all Go tests from one script. Can take ~10 seconds
go test ./auth/
go test ./config/
go test ./misc/
go test ./models/testHelpers/
go test ./models/brand/
go test ./models/tag/
go test ./models/purchase/
go test ./models/user/
