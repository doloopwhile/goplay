require 'rake'

version = -> {
  v = `git describe --tags --always --dirty`.strip
  fail unless $? == 0
  v
}

build_flags = -> {
  %Q{-ldflags "-X main.version \"#{version.()}\""}
}

verbose = '-v' if ENV['VERBOSE']

task :build => :deps do
  sh %Q{go build #{verbose} #{build_flags.()}}
end

task :install => :deps do
  sh %Q{go build #{verbose} #{build_flags.()}}
end

task :deps do
  sh %Q{go get #{verbose} -d}
end
