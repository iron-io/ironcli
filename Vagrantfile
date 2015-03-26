# -*- mode: ruby -*-
# vi: set ft=ruby :

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#
# Place this Vagrantfile in your src folder and run:
#
#     vagrant up
#
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# modified from https://github.com/nathany/vagrant-gopher

# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

Vagrant.require_version ">= 1.5.0"

# See http://dl.golang.org/dl/
GO_ARCHIVES = {
  "linux" => "go1.4.1.linux-amd64.tar.gz",
  "bsd" => "go1.4.1.freebsd-amd64.tar.gz"
}

INSTALL = {
  "linux" => "apt-get update -qq; apt-get install -qq -y git mercurial bzr curl",
  "bsd" => "pkg_add -r git"
}

# location of the Vagrantfile
def gopath
  ENV['GOPATH']
end

# shell script to bootstrap Go
def bootstrap(box)
  install = INSTALL[box]
  archive = GO_ARCHIVES[box]

  profile = <<-PROFILE
  export GOPATH=$HOME/go
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  export CDPATH=.:$GOPATH/src/github.com:$GOPATH/src/code.google.com/p:$GOPATH/src/bitbucket.org:$GOPATH/src/launchpad.net
  PROFILE

  # g++ installation stuff taken from https://gist.github.com/omnus/6404505
  <<-SCRIPT
  #{install}
  if ! [ -f /home/vagrant/#{archive} ]; then
    response=$(curl -O# https://storage.googleapis.com/golang/#{archive})
  fi
  tar -C /usr/local -xzf #{archive}
  echo '#{profile}' >> /home/vagrant/.profile
  sudo add-apt-repository -y ppa:ubuntu-toolchain-r/test
  sudo apt-get -y update
  sudo apt-get -y install g++-4.8
  sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-4.8 50
  echo "\nRun: vagrant ssh #{box} -c 'cd project/path; go test ./...'"
  SCRIPT
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "linux" do |linux|
    linux.vm.box = "ubuntu/trusty64"
    linux.vm.synced_folder gopath, "/home/vagrant/go"
    linux.vm.provision :shell, :inline => bootstrap("linux")
  end
end
