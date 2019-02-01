require 'fileutils'

Given(/^I have "([^"]*)" command installed$/) do |command|
  is_present = system("which #{ command} > /dev/null 2>&1")
  raise "Command #{command} is not present in the system" if not is_present
end

Given("a build of gitsweeper") do

end

Given(/I clone "([^"]*)" repo/) do |repo_name|
  run_silent %(git clone #{repo_name})
end
