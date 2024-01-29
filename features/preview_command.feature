Feature: Preview Command

  Background:
    Given I have "go" command installed
    And a build of gitsweeper

  Scenario: In a Git repo with branches
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test preview`
    Then the output should contain:
      """
      Fetching from the remote...

      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2

      To delete them, run again with `gitsweeper cleanup`
      """
    And the exit status should be 0

  Scenario: In a Git repo with branches and debug enabled
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test --debug preview`
    Then the output should contain "Branch origin/duplicate-branch-1 head (605999f514798915490a1887aa255ea56393de07) was found in master, so has been merged!"
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
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test preview --origin=notexist`
    Then the output should match /Could not find the remote named notexist/
    And the exit status should be 1

  Scenario: Specifying a remote with multiple remotes
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    And I add a new remote "new_remote" with url "http://localhost:8008/dummy-repo.git"
    When I run `gitsweeper-int-test preview --origin=new_remote`
    Then the output should contain:
      """
      Fetching from the remote...

      These branches have been merged into master:
        new_remote/duplicate-branch-1
        new_remote/duplicate-branch-2

      To delete them, run again with `gitsweeper cleanup`
      """
    And the exit status should be 0
