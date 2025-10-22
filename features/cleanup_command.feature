Feature: Cleanup Command

  Background:
    Given I have "go" command installed
    And a build of gitsweeper
    And I have "docker" command installed
    And nothings running on port "8008"

  Scenario: In a git repo with branches with force
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup --force`
    Then the output should contain:
      """
      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2
      
        deleting origin/duplicate-branch-1 - (done)
        deleting origin/duplicate-branch-2 - (done)
      """
    And the exit status should be 0

  Scenario: In a git repo with branches with prompt yes
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup` interactively
    And I type "y"
    Then the output should contain:
      """
      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2
      Delete these branches? [y/n]: 
        deleting origin/duplicate-branch-1 - (done)
        deleting origin/duplicate-branch-2 - (done)
      """
    And the exit status should be 0
  
  Scenario: In a git repo with branches with prompt no
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup` interactively
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
    When I run `gitsweeper-int-test cleanup`
    Then the output should contain:
      """
      Error: This is not a Git repository
      """
    And the exit status should be 1

  Scenario: Specifying a remote with multiple remotes
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    And I add a new remote "new_remote" with url "http://localhost:8008/dummy-repo.git"
    When I run `gitsweeper-int-test cleanup --origin=new_remote --force`
    Then the output should contain:
      """
      Fetching from the remote...

      These branches have been merged into master:
        new_remote/duplicate-branch-1
        new_remote/duplicate-branch-2

        deleting new_remote/duplicate-branch-1 - (done)
        deleting new_remote/duplicate-branch-2 - (done)
      """
  
  Scenario: Specifying a single skip string with debug
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup --force --skip=duplicate-branch-1 --debug`
    Then the output should contain:
      """
      Branch 'origin/duplicate-branch-1' matches skip branch string '[duplicate-branch-1]'
      """
    And the exit status should be 0

  Scenario: Specifying skipping all branches gives specific error message
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup --force --skip=duplicate-branch-1,duplicate-branch-2`
    Then the output should match /No remote branches are available for cleaning up/
    And the exit status should be 0
  
  Scenario: Specifying a non-existant single skip string with debug
    Given no old "gitdocker" containers exist
    And I have a dummy git server running called "gitdocker" running on port "8008"
    And I clone "http://localhost:8008/dummy-repo.git" repo
    And I cd to "dummy-repo"
    When I run `gitsweeper-int-test cleanup --force --skip=skipfakebranch --debug`
    Then the output should contain:
      """
      Fetching from the remote...

      These branches have been merged into master:
        origin/duplicate-branch-1
        origin/duplicate-branch-2

        deleting origin/duplicate-branch-1 - (done)
        deleting origin/duplicate-branch-2 - (done)
      """
    And the exit status should be 0

