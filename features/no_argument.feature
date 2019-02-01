Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of gitsweeper
    When I run `bin/gitsweeper-int-test`
    Then the output should contain:
      """"
      usage: gitsweeper [<flags>] <command> [<args> ...]

      A command-line tool for cleaning up merged branches.
      """"
