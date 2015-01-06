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
  command '/opt/go/bin/goop install'
  environment ({'HOME' => "/home/vagrant",
                "PATH" => "/opt/go/bin:/usr/local/go/bin:#{ENV["PATH"]}"})
  cwd '/home/vagrant/streakers'
end

execute 'copy dotenv.sample to .env' do
  command 'cp -n dotenv.sample .env'
  cwd '/home/vagrant/streakers'
end
