# Cookbook Name:: filesync
# Recipe:: development
#
# Things needed for running filesync

# Other projects (perceptive_idp) default to user "ubuntu", but that user
# doesn't exist on this machine.
node.set['user'] = 'vagrant'

include_recipe 'filesync::default'

include_recipe 'database::postgresql'

default_connection = {username: 'postgres', host: '127.0.0.1'}

postgresql_database_user 'filesync' do
  password 'filesync'
  connection default_connection
  action :create
end

[
  'filesync_development',
  'filesync_test'
].each do |database|
  postgresql_database database do
    owner 'filesync'
    connection default_connection
    action :create
  end

  postgresql_database_user 'filesync' do
    database_name database
    privileges [:all]
    connection default_connection
    action :grant
  end
end

# capture error ouput
# swallow errors if they are because the user exists
# all other errors make the grep also fail
file '/tmp/createuser.err' do
  action :delete
end
execute "sudo -u postgres createuser --superuser #{node[:user]} " +
        '2>/tmp/createuser.err ' +
        '|| grep exists /tmp/createuser.err'

execute 'gem install foreman'

# file sync
execute 'cp dotenv.sample .env' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/filesync"
  not_if { File.file?("/home/#{node[:user]}/filesync/.env") }
end
execute 'bin/bundle' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/filesync"
end
execute 'bin/rake db:create db:migrate' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/filesync"
  not_if { File.file?("/home/#{node[:user]}/filesync/db/schema.rb") }
end
execute 'bin/rake db:reset' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/filesync"
  only_if { File.file?("/home/#{node[:user]}/filesync/db/schema.rb") }
end

include_recipe 'filesync::javascript_runtime'

include_recipe 'perceptive_idp::development'
execute 'bundle install' do
  cwd "/home/#{node[:user]}/perceptive_idp"
end

execute 'echo "Doorkeeper::Application.create!(name: \'File Sync UI\', uid: \'ece6ed48323d4e9191fccf7a45cb7a01ba3311f70f22d791260c3da524a4d9b1\', secret: \'1bb05da868926f5385cfbb62c2184e8cec2b938e80b84e676cdca7cff3fb9b43\', redirect_uri: \'http://localhost:3002/callback.htm\')" | ./bin/rails console' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/perceptive_idp"
end

execute 'echo "Doorkeeper::Application.create!(name: \'File Sync Service\', uid: \'cc79952d8c6cd9a529a3b9f9716be4cc16e75f40124e11a94d6e723acf1248f9\', secret: \'11cd18b0892a2101c64a3a19a742a05b94f6118b43ad5ab87f3d243b4d8a6ea8\', redirect_uri: \'http://localhost:3000\')" | ./bin/rails console' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/perceptive_idp"
end

execute 'echo "User.create!(email: \'admin@example.com\', password: \'password\',  organization: Organization.first)" | ./bin/rails console' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/perceptive_idp"
end

include_recipe 'filesync::webclient'
