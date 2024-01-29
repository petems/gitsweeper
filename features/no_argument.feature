Feature: Version Command

  Background:
    Given I have "go" command installed
    And a build of gitsweeper
    Then the build should be present

  Scenario:
    Given a build of gitsweeper
    When I run `gitsweeper-int-test`
    Then the output should contain:
      """"
      usage: gitsweeper [<flags>] <command> [<args> ...]

      A command-line tool for cleaning up merged branches.
      """"
