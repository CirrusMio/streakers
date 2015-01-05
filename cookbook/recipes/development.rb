# Cookbook Name:: streakers
# Recipe:: development

node.set['go']['owner'] = 'vagrant'

magic_shell_environment 'GOROOT' do
  value '/usr/local/go'
end

magic_shell_environment 'GOPATH' do
  value '/home/vagrant'
end

magic_shell_environment 'GOBIN' do
  value '$GOPATH/bin'
end

include_recipe 'streakers::default'

# TODO(chase): install goop: go get github.com/nitrous-io/goop
# TODO(chase): run `goop install` in ~/streakers
