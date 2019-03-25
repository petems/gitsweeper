Feature: Preview Command

  Background:
    Given I have "go" command installed
    When I run `go build -o bin/gitsweeper ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of gitsweeper
    And a clone of "github.com/petems/gitsweeper"
    And I cd to "gitsweeper"
    When I run `../bin/gitsweeper preview`
    Then the output should match /^Version \d+\.\d\.\d$/