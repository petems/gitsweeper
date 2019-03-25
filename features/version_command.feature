Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o bin/gitsweeper ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of gitsweeper
    When I run `bin/gitsweeper version`
    Then the output should match /^Version \d+\.\d\.\d$/