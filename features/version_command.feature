Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of gitsweeper
    When I run `bin/gitsweeper-int-test version`
    Then the output should contain exactly:
      """"

      0.1.0 development
      """"
