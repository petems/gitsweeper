require 'aruba/cucumber'
require 'docker'
require 'fileutils'
require 'forwardable'
require 'tmpdir'

$bin_dir = File.expand_path('../../../bin/', __FILE__)
$aruba_dir = File.expand_path('../../..', __FILE__) + '/tmp/aruba'

Aruba.configure do |config|
  # increase process exit timeout from the default of 3 seconds
  config.exit_timeout = 20
  # allow absolute paths for tests involving no repo
  config.allow_absolute_paths = true
  # don't be "helpful"
  config.remove_ansi_escape_sequences = false
end

Before do
  aruba.environment.update(
    'PATH' => "#{$bin_dir}:#{ENV['PATH']}",
  )
  FileUtils.rm_rf("#{$aruba_dir}/bare-git-repo")
  FileUtils.rm_rf("#{$aruba_dir}/dummy-repo")
  FileUtils.rm_rf("#{$bin_dir}/gitsweeper-int-test")
end

After do
  begin
    app = Docker::Container.get("gitdocker")
    app.delete(force: true)
  rescue Docker::Error::NotFoundError
  end
  FileUtils.rm_rf("#{$aruba_dir}/bare-git-repo")
  FileUtils.rm_rf("#{$aruba_dir}/dummy-repo")
  FileUtils.rm_rf("#{$bin_dir}/gitsweeper-int-test")
end
