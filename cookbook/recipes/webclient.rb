# Cookbook Name:: filesync
# Recipe:: webclient
#
# Copyright (c) 2014 Lexmark International Technology S.A.  All rights reserved.
# All rights reserved - Do Not Redistribute
#
node.set['user'] = 'vagrant'

execute 'npm update' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/fns-evolution-js-webclient"
end

execute 'npm install' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/fns-evolution-js-webclient"
end

execute 'sudo npm install -g grunt-cli' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/fns-evolution-js-webclient"
end

execute 'sudo npm install -g http-server' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/fns-evolution-js-webclient"
end

execute 'grunt dev --brand=recall.com' do
  user node[:user]
  environment 'HOME' => "/home/#{node[:user]}"
  cwd "/home/#{node[:user]}/fns-evolution-js-webclient"
end
