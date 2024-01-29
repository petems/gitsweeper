Feature: Preview Command

  Background:
    Given I have "go" command installed
    And a build of gitsweeper

  Scenario: In a Git repo with branches
    Given I clone "https://github.com/petems/example-repo-with-remote-branches.git" repo
    And I cd to "example-repo-with-remote-branches"
    When I run `gitsweeper-int-test preview`
    Then the output should contain:
      """
      Fetching from the remote...

      These branches have been merged into master:
        origin/branch_thats_been_merged_to_master

      To delete them, run again with `gitsweeper cleanup`
      """
    And the exit status should be 0

  Scenario: In a Git repo with branches and debug enabled
    Given I clone "https://github.com/petems/example-repo-with-remote-branches.git" repo
    And I cd to "example-repo-with-remote-branches"
    When I run `gitsweeper-int-test --debug preview`
    Then the output should contain "Branch origin/branch_thats_been_merged_to_master head (574d05288ddb0449402829e90b5752d78a7f54d7) was found in master, so has been merged!"
    And the exit status should be 0

  Scenario: In a non-git repo
    Given I run `mkdir -p not-a-git-repo`
    And I cd to "not-a-git-repo"
    When I run `gitsweeper-int-test preview`
    Then the output should match /This is not a Git repository/
    And the exit status should be 1

  Scenario: In a non-git repo
    Given I create a bare git repo called "bare-git-repo"
    And I cd to "bare-git-repo"
    When I run `gitsweeper-int-test preview`
    Then the output should match /Could not find the remote named origin/
    And the exit status should be 1

  Scenario: Using a non-existant remote
    Given I clone "https://github.com/petems/example-repo-with-remote-branches.git" repo
    And I cd to "example-repo-with-remote-branches"
    When I run `gitsweeper-int-test preview --origin=notexist`
    Then the output should match /Could not find the remote named notexist/
    And the exit status should be 1
