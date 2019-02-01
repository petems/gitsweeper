Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: Version with no flags
    Given a build of gitsweeper
    When I run `bin/gitsweeper-int-test version`
    Then the output should contain exactly:
      """""

      0.1.0 development
      """""

  Scenario: Version with --debug flag
    Given a build of gitsweeper
    When I run `bin/gitsweeper-int-test --debug version`
    Then the output should match /--debug setting detected - Info level logs enabled/
