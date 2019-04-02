Feature: Preview Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: In a Git repo with branches
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

  Scenario: In a non-git repo
    Given I run `mkdir -p not-a-git-repo`
    And I cd to "not-a-git-repo"
    When I run `bin/gitsweeper-int-test preview`
    Then the output should contain:
      """

      gitsweeper-int-test: error: Error when looking for branches repository does not exist
      """
    And the exit status should be 1

  Scenario: In a non-git repo
    Given I create a bare git repo called "bare-git-repo"
    And I cd to "bare-git-repo"
    When I run `bin/gitsweeper-int-test preview`
    Then the output should match /reference not found/
    And the exit status should be 1
