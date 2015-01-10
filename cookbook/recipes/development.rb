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
  command 'cp dotenv.sample .env'
  cwd '/home/vagrant/streakers'
end

node.set['user'] = 'vagrant'

include_recipe 'streakers::default'

# capture error ouput
# swallow errors if they are because the user exists
# all other errors make the grep also fail
file '/tmp/createuser.err' do
  action :delete
end
execute "sudo -u postgres createuser --superuser #{node[:user]} " +
        '2>/tmp/createuser.err ' +
        '|| grep exists /tmp/createuser.err'

include_recipe 'database::postgresql'

default_connection = {username: 'postgres', host: '127.0.0.1'}

postgresql_database_user 'streaker' do
  password 'streaker'
  connection default_connection
  action :create
end

[
  'streaker_development',
  'streaker_test'
].each do |database|
  postgresql_database database do
    owner 'streaker'
    connection default_connection
    action :create
  end

  postgresql_database_user 'streaker' do
    database_name database
    privileges [:all]
    connection default_connection
    action :grant
  end
end
