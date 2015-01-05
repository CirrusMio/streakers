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

execute 'go get github.com/nitrous-io/goop' do
  user 'vagrant'
  environment 'HOME' => '/home/vagrant'
  cwd '/home/vagrant/streakers'
end

execute 'goop install' do
  user 'vagrant'
  environment 'HOME' => '/home/vagrant'
  cwd '/home/vagrant/streakers'
end
