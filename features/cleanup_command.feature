Feature: Cleanup Command

  Background:
    Given I have "go" command installed
    And I have "docker" command installed
    And nothings running on port "8008"
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: In a git repo with branches with force
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `bin/gitsweeper-int-test cleanup --force`
    Then the output should contain:
      """
      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2
      
        deleting duplicate-branch-1 - (done)
        deleting duplicate-branch-2 - (done)
      """
    And the exit status should be 0

  Scenario: In a git repo with branches with prompt yes
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `bin/gitsweeper-int-test cleanup` interactively
    And I type "y"
    Then the output should contain:
      """
      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2
      Delete these branches? [y/n]: 
        deleting duplicate-branch-1 - (done)
        deleting duplicate-branch-2 - (done)
      """
    And the exit status should be 0
  
  Scenario: In a git repo with branches with prompt no
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `bin/gitsweeper-int-test cleanup` interactively
    And I type "n"
    Then the output should contain:
      """
      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2
      Delete these branches? [y/n]: OK, aborting.
      """
    And the exit status should be 0

  Scenario: In a non-git repo
    Given I run `mkdir -p not-a-git-repo`
    And I cd to "not-a-git-repo"
    When I run `bin/gitsweeper-int-test cleanup`
    Then the output should contain:
      """

      gitsweeper-int-test: error: Error when looking for branches repository does not exist
      """
    And the exit status should be 1
