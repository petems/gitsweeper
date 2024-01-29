Feature: Version Command

  Scenario: Version with no flags
    Given a build of gitsweeper
    When I run `gitsweeper-int-test version`
    Then the output should contain exactly:
      """""
      0.1.0 development
      """""

  Scenario: Version with --debug flag
    Given a build of gitsweeper
    When I run `gitsweeper-int-test --debug version`
    Then the output should match /--debug setting detected - Info level logs enabled/
