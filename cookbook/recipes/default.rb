# Cookbook Name:: streakers
# Recipe:: default

# prerequisites
include_recipe 'annoyances'
include_recipe 'build-essential'
include_recipe 'git'
include_recipe 'golang'
include_recipe 'golang::packages'

node.set['postgresql']['password']['postgres'] = 'password'
include_recipe 'postgresql::server'
include_recipe 'postgresql::client'
