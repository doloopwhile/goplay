require 'rake'

version = -> {
  v = `git describe --tags --always --dirty`.strip
  fail unless $? == 0
  v
}

ldflags = -> {
  %Q{"-X main.version \"#{version.()}\""}
}

verbose = '-v' if ENV['VERBOSE']

task :build => :deps do
  sh %Q{go build #{verbose} -ldflags #{ldflags.()}}
end

task :install => :deps do
  sh %Q{go build #{verbose} -ldflags #{ldflags.()}}
end

task :deps do
  sh %Q{go get #{verbose} -d}
end

task :goxc do
  sh %Q{goxc -tasks='xc archive' -d dist -bc='linux,!arm windows,386 darwin' -build-ldflags=#{ldflags.()}}
end
