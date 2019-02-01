Feature: Preview Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given I clone "git://github.com/petems/example-repo-with-remote-branches.git" repo
    And I cd to "example-repo-with-remote-branches"
    When I run `bin/gitsweeper-int-test preview`
    Then the output should contain:
      """

      Fetching from the remote...

      These branches have been merged into master:
        origin/branch_thats_been_merged_to_master
      """
    And the exit status should be 0
