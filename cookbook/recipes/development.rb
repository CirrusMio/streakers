# Cookbook Name:: streakers
# Recipe:: development

include_recipe 'streakers::default'

magic_shell_environment 'GOROOT' do
  value '/usr/local/go'
end

magic_shell_environment 'GOPATH' do
  value '/home/vagrant/streakers'
end

# TODO(chase): install goop: go get github.com/nitrous-io/goop
# TODO(chase): run `goop install` in ~/streakers
