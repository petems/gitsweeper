require 'aruba/cucumber'
require 'fileutils'
require 'forwardable'
require 'tmpdir'

bin_dir = File.expand_path('../fakebin', __FILE__)
system_git = `which git 2>/dev/null`.chomp

Before do
  # increase process exit timeout from the default of 3 seconds
  @aruba_timeout_seconds = 10
  # don't be "helpful"
  @aruba_keep_ansi = true
end

After do
  @server.stop if defined? @server and @server
  FileUtils.rm_f("#{bin_dir}/vim")
end

RSpec::Matchers.define :be_successful_command do
  match do |cmd|
    cmd.success?
  end

  failure_message do |cmd|
    %(command "#{cmd}" exited with status #{cmd.status}:) <<
      cmd.output.gsub(/^/, ' ' * 2)
  end
end

class SimpleCommand
  attr_reader :output
  extend Forwardable

  def_delegator :@status, :exitstatus, :status
  def_delegators :@status, :success?

  def initialize cmd
    @cmd = cmd
  end

  def to_s
    @cmd
  end

  def self.run cmd
    command = new(cmd)
    command.run
    command
  end

  def run
    @output = `#{@cmd} 2>&1`.chomp
    @status = $?
    $?.success?
  end
end

World Module.new {
  # If there are multiple inputs, e.g., type in username and then type in password etc.,
  # the Go program will freeze on the second input. Giving it a small time interval
  # temporarily solves the problem.
  # See https://github.com/cucumber/aruba/blob/7afbc5c0cbae9c9a946d70c4c2735ccb86e00f08/lib/aruba/api.rb#L379-L382
  def type(*args)
    super.tap { sleep 0.1 }
  end

  def history
    histfile = File.join(ENV['HOME'], '.history')
    if File.exist? histfile
      File.readlines histfile
    else
      []
    end
  end

  def assert_command_run cmd
    cmd += "\n" unless cmd[-1..-1] == "\n"
    expect(history).to include(cmd)
  end

  def edit_hub_config
    config = File.join(ENV['HOME'], '.config/hub')
    FileUtils.mkdir_p File.dirname(config)
    if File.exist? config
      data = YAML.load File.read(config)
    else
      data = {}
    end
    yield data
    File.open(config, 'w') { |cfg| cfg << YAML.dump(data) }
  end

  define_method(:text_editor_script) do |bash_code|
    File.open("#{bin_dir}/vim", 'w', 0755) { |exe|
      exe.puts "#!/bin/bash"
      exe.puts "set -e"
      exe.puts bash_code
    }
  end

  def run_silent cmd
    in_current_dir do
      command = SimpleCommand.run(cmd)
      expect(command).to be_successful_command
      command.output
    end
  end

  def empty_commit(message = nil)
    unless message
      @empty_commit_count = defined?(@empty_commit_count) ? @empty_commit_count + 1 : 1
      message = "empty #{@empty_commit_count}"
    end
    run_silent "git commit --quiet -m '#{message}' --allow-empty"
  end

  # Aruba unnecessarily creates new Announcer instance on each invocation
  def announcer
    @announcer ||= super
  end

  def shell_escape(message)
    message.to_s.gsub(/['"\\ $]/) { |m| "\\#{m}" }
  end

  %w[output_from stdout_from stderr_from all_stdout all_stderr].each do |m|
    define_method(m) do |*args|
      home = ENV['HOME'].to_s
      output = super(*args)
      if home.empty?
        output
      else
        output.gsub(home, '$HOME')
      end
    end
  end
}
