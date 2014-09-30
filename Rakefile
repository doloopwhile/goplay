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

task :gox do
  sh %Q{gox -output="dist/{{.Dir}}.{{.OS}}_{{.Arch}}" -arch="386 amd64" -os 'linux windows darwin' -ldflags #{ldflags.()}}
end
