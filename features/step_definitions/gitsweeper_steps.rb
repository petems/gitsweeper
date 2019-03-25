Given(/^I have "([^"]*)" command installed$/) do |command|
  is_present = system("which #{ command} > /dev/null 2>&1")
  raise "Command #{command} is not present in the system" if not is_present
end

Given("a build of gitsweeper") do

end

Given("a clone of github.com/petems/gitsweeper") do 
  
end

Given("a clone of {string}") do |string|
  is_present = system(%(git checkout --quiet #{string}))
end