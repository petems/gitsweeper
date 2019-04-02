Feature: Cleanup Command

  Background:
    Given I have "go" command installed
    And I have "docker" command installed
    And nothings running on port "8008"
    When I run `go build -o ../../bin/gitsweeper-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: In a git repo with branches
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `bin/gitsweeper-int-test cleanup`
    Then the output should contain:
      """

      Fetching from the remote...
      deleting duplicate-branch-1 - (done)
      deleting duplicate-branch-2 - (done)
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
