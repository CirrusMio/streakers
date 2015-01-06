# Cookbook Name:: streakers
# Recipe:: development

node.set['go']['version'] = '1.4'
node.set['go']['owner'] = 'vagrant'
node.set['go']['packages'] = ['github.com/nitrous-io/goop']

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

execute 'install go packages' do
  command 'goop install'
  cwd '/home/vagrant/streakers'
end
