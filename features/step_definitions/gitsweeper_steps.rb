require 'fileutils'

Given(/^I have "([^"]*)" command installed$/) do |command|
  is_present = system("which #{ command} > /dev/null 2>&1")
  raise "Command #{command} is not present in the system" if not is_present
end

Given("a build of gitsweeper") do

end

Given(/nothings running on port "(\w+)"/) do |port|
  running_on_port = system("lsof -i TCP:#{port} > /dev/null 2>&1")
  raise "Something running on port #{port}" if running_on_port
end

Given /^no old "(\w+)" containers exist$/ do |container_name|
  begin
    app = Docker::Container.get(container_name)
    app.delete(force: true)
  rescue Docker::Error::NotFoundError
  end
end

Given /^I have a dummy git server running called "(\w+)" running on port "(\w+)"$/ do |container_name, port|
  steps %Q(
    Given no old "#{container_name}" containers exist
    When I successfully run `docker run -d -p '#{port}:80' --name='#{container_name}' petems/dummy-git-repo`
  )
  sleep 3
end

Given(/I clone "([^"]*)" repo/) do |repo_name|
  run_silent %(git clone #{repo_name})
end

Given(/I create a bare git repo called "([^"]*)"/) do |repo_name|
  run_silent %(git init --bare #{repo_name})
end
