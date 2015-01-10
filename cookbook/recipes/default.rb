# Cookbook Name:: streakers
# Recipe:: default

# prerequisites
include_recipe 'annoyances'
include_recipe 'build-essential'
include_recipe 'git'
include_recipe 'golang'
include_recipe 'golang::packages'

package('libcurl4-openssl-dev').run_action(:install)
package('libxml2-dev').run_action(:install)
package('libxslt1-dev').run_action(:install)
package('libpq-dev').run_action(:install)

node.set['postgresql']['password']['postgres'] = 'password'

include_recipe 'postgresql::server'
include_recipe 'postgresql::client'
